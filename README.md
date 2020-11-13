## jx-git-operator

`jx-git-operator` is an operator which polls a git repository for changes and triggers a Kubernetes `Job` to process changes in git.

It can be used to install/upgrade any environment (development, staging, production) via some GitOps approach using some set of tools (`kubectl`, `helm`, `helmfile`, `kpt`, `kustomize` etc).


### Overview

This operator will poll for git commits in git repositories. If a new git commit sha is detected in the repository, the repository is cloned and a `Job` is triggered.

The `Job` to trigger is defined in the `.jx/git-operator/job.yaml` file inside the git repository.

Here is [an example repository](https://github.com/jenkins-x/jx3-boot-config/tree/master/.jx/git-operator) so you can see how it works

### Installing the easy way

The `jx-admin` command line will install the operator for you and setup the Git URL Secret so try one of these commands:

* [jx-admin operator](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_operator.md) if you already have a git repository and just want to install the operator


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

To see the logs of the operator try:


```bash
kubectl logs -f -l app=jx-git-operator
```    

you should see it polling your git repository and triggering `Job` instances whenever a change is deteted


### Running 

You can run the `jx-git-operator` locally on the command line if you want. Actions will be created as Kubernetes Jobs even if you run the binary locally - it is just the git polling which runs locally.

Download the [jx-git-operator binary](https://github.com/jenkins-x/jx-git-operator/releases) for your operating system and add it to your `$PATH`.
