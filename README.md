## jx-git-operator

`jx-git-operator` is an operator which polls a git repository for changes and triggers a Kubernetes `Job` to process changes in git.

The definition of the `Job` is defined in the git repository leaving you free to trigger any kind of `Job` you like. e.g. use `kubectl apply` if you wish or `helm install` or `kustomize` or whatever. 

The `jx-git-operator` is small with a minimal footprint and has no dependencies so can be used to install/upgrade/configure anything you like in any cluster.

It can be used to install/upgrade any environment (development, staging, production) via a GitOps approach using any set of tools you like ([helm](https://helm.sh/), [helmfile](https://github.com/roboll/helmfile), [jx](https://github.com/jenkins-x/jx-cli/releases),  [kapp](https://get-kapp.io/), [kpt](https://googlecontainertools.github.io/kpt/), [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/), [kustomize](https://kustomize.io/) etc).

You are in full control over exactly what the `Job` does in each cluster.   

### How it works

The `jx-git-operator` will poll for git commits in git repositories. If a new git commit sha is detected in the repository, the repository is cloned and a `Job` is triggered for that sha.

The `Job` to trigger is defined in the git repository being polled. The default file is looked for at `versionStream/git-operator/job.yaml` or `.jx/git-operator/job.yaml`.

Here is [an example repository](https://github.com/jx3-gitops-repositories/jx3-kubernetes) with a [versionStream/git-operator/job.yaml](https://github.com/jx3-gitops-repositories/jx3-kubernetes/blob/master/versionStream/git-operator/job.yaml) so you can see how it works

The Jenkins X `Job` uses a simple `Makefile` to trigger steps in the git operator job making it super easy for you to use any permutation of commands using tools like ([helm](https://helm.sh/), [helmfile](https://github.com/roboll/helmfile), [jx](https://github.com/jenkins-x/jx-cli/releases), [kapp](https://get-kapp.io/), [kpt](https://googlecontainertools.github.io/kpt/), [kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/), [kustomize](https://kustomize.io/) which are also trivial to test locally via running `make` in a local git clone of the git repository.

This lets you define the exact GitOps process you wish to use without being locked into a specific operators decisions on what is supported.
                 

### Installing the easy way

The [jx admin operator](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_operator.md) command line will install the operator for you and tail the log of the triggered `Job` so you can see what its doing. 

Under the covers the command will download the [helm](https://helm.sh/) binary for your platform, output the helm command to install the operator then actually run the command for you.

See the [Getting Started Instructions](https://jenkins-x.io/v3/admin/setup/operator/)

Note that if you are using [terraform](https://www.terraform.io/) using a Jenkins X terraform module then the git operator is automatically installed into the Kubernetes clusters via Terraform.
      

### Installing the hard way

To install the git operator by hand using [helm 3](https://helm.sh/) then try:

Setup a namespace:

```bash 
helm repo add jx3 https://storage.googleapis.com/jenkinsxio/charts
helm upgrade --install \
    --set url=$GIT_URL \
    --set username=$GIT_USER \
    --set password=$GIT_TOKEN \
     jx-git-operator --create-namespace jxgo jx3/jx-git-operator
```

### Viewing the logs
      

#### kubectl 

To see the new logs of the operator try:


```bash
kubectl logs -f -l app=jx-git-operator -n jx-git-operator
```    

you should see it polling your git repository and triggering `Job` instances whenever a change is detected

#### jx

To view the logs of the jobs triggered by the git operator you can use the [jx admin log](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_log.md) command:

```bash 
jx admin log
```
         
This commmand will list all the known `Job` instances sorted in time order, letting you pick one then showing the log details.

If you know you have just done a git commit and are waiting for the boot job to start you can run:

```bash 
jx admin log --wait
```

Which will wait for a running `Job` to display.


#### Octant

If you use the [Jenkins X plugin](https://github.com/jenkins-x/octant-jx) for [Octant](https://octant.dev/) via:

```bash 
jx ui
```
                                               
Then you can view the boot Jobs triggered by the git operator (along with the commit message, user and timestamp of the git commits) in the Boot Jobs tab.

See the [Jenkins X Console documentation for more](https://jenkins-x.io/v3/develop/ui/octant/)


### Running locally

You can run the `jx-git-operator` locally on the command line if you want. Actions will be created as Kubernetes Jobs even if you run the binary locally - it is just the git polling which runs locally.

Download the [jx-git-operator binary](https://github.com/jenkins-x/jx-git-operator/releases) for your operating system and add it to your `$PATH`.
