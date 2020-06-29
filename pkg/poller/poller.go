package poller

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher/job"
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/jenkins-x/jx-git-operator/pkg/repo/secret"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/files"
	"github.com/jenkins-x/jx-helpers/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/pkg/gitclient/cli"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Options the configuration options for the poller
type Options struct {
	GitClient  gitclient.Interface
	RepoClient repo.Interface
	Launcher   launcher.Interface

	// CommandRunner used to run git commands if no GitClient provided
	CommandRunner cmdrunner.CommandRunner

	// KubeClient is used to lazy create the repo client and launcher
	KubeClient kubernetes.Interface

	Dir             string        `env:"WORK_DIR"`
	Namespace       string        `env:"NAMESPACE"`
	GitBinary       string        `env:"GIT_BINARY"`
	PollDuration    time.Duration `env:"POLL_DURATION"`
	NoLoop          bool          `env:"NO_LOOP"`
	NoResourceApply bool          `env:"NO_RESOURCE_APPLY"`
}

// Run polls for git changes
func (o *Options) Run() error {
	err := o.validateOptions()
	if err != nil {
		return errors.Wrap(err, "invalid options")
	}

	if o.Namespace != "" {
		log.Logger().Infof("looking in namespace %s for Secret resources with selector %s", o.Namespace, constants.DefaultSelector)
	}

	if !o.NoLoop {
		log.Logger().Infof("using poll duration %s", o.PollDuration.String())
	}
	for {
		err = o.Poll()
		if err != nil {
			return err
		}
		if o.NoLoop {
			return nil
		}
		time.Sleep(o.PollDuration)
	}
}

// Poll polls the available repositories
func (o *Options) Poll() error {
	err := o.validateOptions()
	if err != nil {
		return errors.Wrap(err, "invalid options")
	}

	repos, err := o.RepoClient.List()
	if err != nil {
		return errors.Wrapf(err, "failed to list repositories")
	}

	if len(repos) == 0 {
		log.Logger().Infof("no repositories found")
		return nil
	}
	for _, r := range repos {
		err = o.pollRepository(r)
		if err != nil {
			return errors.Wrapf(err, "failed to poll repository %s in namespace %s", r.Name, r.Namespace)
		}
	}
	return nil
}

func (o *Options) pollRepository(r repo.Repository) error {
	name := r.Name
	log.Logger().Infof("polling repository %s in namespace %s with git URL %s", name, r.Namespace, r.GitURL)

	dir := filepath.Join(o.Dir, name)
	exists, err := files.DirExists(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to check dir exists %s", dir)
	}
	if !exists {
		log.Logger().Infof("cloning repository %s to %s", name, dir)
		_, err = o.GitClient.Command(o.Dir, "clone", r.GitURL, dir)
		if err != nil {
			return errors.Wrapf(err, "failed to clone repository %s", name)
		}
	} else {
		_, err = o.GitClient.Command(dir, "pull", "origin", "master")
		if err != nil {
			return errors.Wrapf(err, "failed to pull repository %s", name)
		}
	}
	text, err := o.GitClient.Command(dir, "rev-parse", "HEAD")
	if err != nil {
		return errors.Wrapf(err, "failed to find latest commit sha for repository %s", name)
	}
	text = strings.TrimSpace(text)
	log.Logger().Infof("repository %s has latest commit sha %s", name, text)

	if text == "" {
		return errors.Errorf("could not find latest commit sha for repository %s", name)
	}

	_, err = o.Launcher.Launch(launcher.LaunchOptions{
		Repository:      r,
		GitSHA:          text,
		Dir:             dir,
		NoResourceApply: o.NoResourceApply,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to launch job for %s", name)
	}
	return nil
}

// validateOptions validates the options and lazily creates any resources required
func (o *Options) validateOptions() error {
	if o.PollDuration.Milliseconds() == int64(0) {
		o.PollDuration = time.Second * 30
	}
	if o.GitClient == nil {
		o.GitClient = cli.NewCLIClient(o.GitBinary, o.CommandRunner)
	}
	var err error
	if o.RepoClient == nil {
		o.RepoClient, err = secret.NewClient(o.KubeClient, o.Namespace, constants.DefaultSelector)
		if err != nil {
			return errors.Wrapf(err, "failed to create repo client")
		}
	}
	if o.Launcher == nil {
		o.Launcher, err = job.NewLauncher(o.KubeClient, o.Namespace, constants.DefaultSelector, o.CommandRunner)
		if err != nil {
			return errors.Wrapf(err, "failed to create launcher")
		}
	}
	if o.Dir == "" {
		o.Dir, err = ioutil.TempDir("", "jx-git-operator-")
		if err != nil {
			return errors.Wrapf(err, "failed to create temp dir")
		}
	}
	return nil
}
