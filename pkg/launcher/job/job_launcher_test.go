package job_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jenkins-x/jx-helpers/v3/pkg/yamls"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher/job"
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestJobLauncher(t *testing.T) {
	ns := "jx"
	repoName := "fake-repository"
	gitURL := "https://fakeuser:fakepwd@github.com/jenkins-x/fake-repository.git"
	gitSha := "dummysha1234"
	lastCommitAuthor := "jstrachan"
	lastCommitAuthorEmail := "something@gmail.com"
	lastCommitDate := "Wed, 24 Feb 2021 10:13:14 +0000"
	lastCommitMessage := "fix: upgrading my app"
	repoURL := "https://github.com/jenkins-x/fake-repository.git"
	lastCommitURL := strings.TrimSuffix(repoURL, ".git") + "/commit/" + gitSha

	fs, err := os.ReadDir("test_data")
	require.NoError(t, err, "failed to load test data")

	for _, f := range fs {
		if f == nil || !f.IsDir() {
			continue
		}
		name := f.Name()

		t.Logf("running test %s\n", name)
		resourcesDir, err := filepath.Abs(filepath.Join("test_data", name, "versionStream", "git-operator", "resources"))
		require.NoError(t, err, "failed to get absolute dir %s", resourcesDir)

		kubeClient := fake.NewSimpleClientset(
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      repoName,
					Namespace: ns,
					Labels: map[string]string{
						constants.DefaultSelectorKey: constants.DefaultSelectorValue,
					},
				},
				Data: map[string][]byte{
					"url": []byte(gitURL),
				},
			},
		)
		runner := &fakerunner.FakeRunner{}

		client, err := job.NewLauncher(kubeClient, ns, constants.DefaultSelector, runner.Run)
		require.NoError(t, err, "failed to create launcher client")

		o := &launcher.LaunchOptions{
			Repository: repo.Repository{
				Name:      repoName,
				Namespace: ns,
				GitURL:    gitURL,
			},
			GitSHA:                gitSha,
			LastCommitAuthor:      lastCommitAuthor,
			LastCommitAuthorEmail: lastCommitAuthorEmail,
			LastCommitDate:        lastCommitDate,
			LastCommitMessage:     lastCommitMessage,
			Dir:                   filepath.Join("test_data", name),
		}
		objects, err := client.Launch(o)
		require.NoError(t, err, "failed to launch the job")
		require.Len(t, objects, 1, "should have created one runtime.Object after launching")

		o1 := objects[0]
		j1, ok := o1.(*v1.Job)
		require.True(t, ok, "could not convert object %#v to a Job")

		t.Logf("created Job with name %s", j1.Name)

		msg := "created Job"
		testhelpers.AssertLabel(t, constants.DefaultSelectorKey, constants.DefaultSelectorValue, j1.ObjectMeta, msg)
		testhelpers.AssertLabel(t, launcher.RepositoryLabelKey, repoName, j1.ObjectMeta, msg)
		testhelpers.AssertLabel(t, launcher.CommitShaLabelKey, gitSha, j1.ObjectMeta, msg)
		testhelpers.AssertAnnotation(t, launcher.CommitAuthorAnnotation, lastCommitAuthor, j1.ObjectMeta, msg)
		testhelpers.AssertAnnotation(t, launcher.CommitAuthorEmailAnnotation, lastCommitAuthorEmail, j1.ObjectMeta, msg)
		testhelpers.AssertAnnotation(t, launcher.CommitDateAnnotation, lastCommitDate, j1.ObjectMeta, msg)
		testhelpers.AssertAnnotation(t, launcher.CommitMessageAnnotation, lastCommitMessage, j1.ObjectMeta, msg)
		testhelpers.AssertAnnotation(t, launcher.CommitURLAnnotation, lastCommitURL, j1.ObjectMeta, msg)
		testhelpers.AssertAnnotation(t, launcher.RepositoryURLAnnotation, repoURL, j1.ObjectMeta, msg)

		runner.ExpectResults(t,
			fakerunner.FakeResult{
				CLI: "kubectl apply -f " + resourcesDir,
			},
		)

		// we should not recreate the Job if we try to launch again as it already exists
		objects, err = client.Launch(o)
		require.NoError(t, err, "failed to launch the job")
		require.Len(t, objects, 0, "should not have a created a runtime.Object as we already have one for the commit sha")

		if name == "customjob" {
			containers := j1.Spec.Template.Spec.Containers
			require.Len(t, containers, 2, "containers for test %s", name)

			c2 := containers[1]
			assert.Equal(t, "gsm", c2.Name, "container[1].Name for test %s", name)
			t.Logf("generated gsm sidecar")
		}
	}
}

func TestOverlayJob(t *testing.T) {
	jsonMap := &strategicpatch.JSONMap{}

	path := filepath.Join("test_data", "somerepo", "versionStream", "git-operator", "job.yaml")
	err := yamls.LoadFile(path, jsonMap)
	require.NoError(t, err, "failed to load file %s", path)

	overlay := &v1.Job{
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "job",
							Env: []corev1.EnvVar{
								{
									Name:  "SOME_NAME",
									Value: "SOME_NAME_NEW_VALUE",
								},
								{
									Name:  "MY_NEW_ENV",
									Value: "MY_NEW_ENV_VALUE",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{Name: "oci-secret", MountPath: "/root/.docker"},
							},
						},
					},
					Volumes: []corev1.Volume{{
						Name: "oci-secret",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "oci",
							},
						}},
					},
				},
			},
		},
	}

	converter := runtime.DefaultUnstructuredConverter
	overLayMap, err := converter.ToUnstructured(&overlay)
	require.NoError(t, err, "failed to convert overlay to map")

	enriched, err := job.OverlayJob(jsonMap, (*strategicpatch.JSONMap)(&overLayMap))

	require.NoError(t, err, "failed to apply overlay")

	vsJob := &v1.Job{}
	err = converter.FromUnstructuredWithValidation(enriched, vsJob, true)
	require.NoError(t, err, "failed to convert unstructured to job")

	containers := vsJob.Spec.Template.Spec.Containers
	require.Len(t, containers, 1, "job should have 1 container")
	container := &containers[0]
	require.Equal(t, "job", container.Name, "container[0].Name")

	env := container.Env

	AssertEnvValue(t, container, "SOME_NAME", "SOME_NAME_NEW_VALUE", "job.container[0]")
	AssertEnvValue(t, container, "MY_NEW_ENV", "MY_NEW_ENV_VALUE", "job.container[0]")
	AssertVolumeMount(t, container, "oci-secret", "/root/.docker")
	AssertVolume(t, &vsJob.Spec.Template.Spec, "oci-secret", "oci")

	require.Len(t, env, 2, "container[0].Env")

}

func AssertVolumeMount(t *testing.T, container *corev1.Container, volumeName, expectedPath string) {
	for _, vm := range container.VolumeMounts {
		if vm.Name == volumeName {
			assert.Equal(t, expectedPath, vm.MountPath)
			return
		}
	}
	assert.Fail(t, "missing volume mount %s in container %s", volumeName, container.Name)
}

func AssertVolume(t *testing.T, podSpec *corev1.PodSpec, volumeName, expectedSecret string) {
	//nolint:gocritic
	for _, vm := range podSpec.Volumes {
		if vm.Name == volumeName {
			assert.Equal(t, expectedSecret, vm.VolumeSource.Secret.SecretName)
			return
		}
	}
	assert.Fail(t, "missing volume %s", volumeName)
}

func AssertEnvValue(t *testing.T, container *corev1.Container, envName, expectedValue, message string) {
	for _, e := range container.Env {
		if e.Name == envName {
			assert.Equal(t, expectedValue, e.Value, "envVar %s in container: %s for %s", envName, container.Name, message)
			return
		}
	}
	assert.Fail(t, "missing envVar %s in container: %s for %s", envName, container.Name, message)
}
