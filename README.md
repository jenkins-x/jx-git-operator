## jx-git-operator

`jx-git-operator` is an operator which polls a git repository for changes and triggers a Kubernetes `Job` to process changes in git.

It can be used to install/upgrade any environment (development, staging, production) via some GitOps approach using some set of tools (`kubectl`, `helm`, `helmfile`, `kpt`, `kustomize` etc).


### Overview

This operator will poll for git commits in git repositories. If a new git commit sha is detected in the repository, the repository is cloned and a `Job` is triggered.

The `Job` to trigger is defined in the `.jx/git-operator/job.yaml` file inside the git repository.

Here is [an example repository](https://github.com/jenkins-x/jx3-boot-config/tree/master/.jx/git-operator) so you can see how it works

### Installing the easy way

The `jx-admin` command line will install the operator for you and setup the Git URL Secret so try one of these commands:

* [jx-admin create](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_create.md) if you don't yet have a git repository 
* [jx-admin operator](https://github.com/jenkins-x/jx-admin/blob/master/docs/cmd/jx-admin_operator.md) if you already have a git repository and just want to install the operator


### Installing the hard way

To install the git operator by hand using [helm 3](https://helm.sh/) then try:

Setup a namespace:

```bash 
kubectl create ns jx-git-operator
jx ns jx-git-operator
```

Then use helm to install/upgrade:
         
```bash    
helm repo add jx-labs https://storage.googleapis.com/jenkinsxio-labs-private/charts
helm install jxgo jx-labs/jx-git-operator
```

You can configure the git polling frequency via the `env.POLL_DURATION` property which supports go `time.Duration` syntax such as `10m` or `40s`

The chart defaults to using a `cluster-admin` role so it can create `Job` resources in any namespace along with any associated resources specified in a git repository at `.jx/git-operator/resources/*.yaml`

You can enable strict mode which only requires roles to read `Secret` resources in the namespace its installed and list/create `Job` resources via the `rbac.strict = true`. 

To avoid cluster roles use `rbac.cluster = false` which only uses a `Role` and `RoleBinding` in current namespace.

### Setting up a repository

The git repository you wish to boot needs to have the `.jx/git-operator/job.yaml` defined to specify the Kubernetes `Job` to perform the boot job.

A `Job` needs to have an associated `ServiceAccount` and either a `ClusterRole` + `ClusterRoleBinding` or `Role` + `RoleBinding`. You can specify those additional resources in the `.jx/git-operator/resources/*.yaml` directory and the operator will `kubectl apply -f .jx/git-operator/resources` before creating the `Job`.

You can disable this behavior by using `rbac.strict = true` when installing the operator. In this case an administrator will need to run: `kubectl apply -f .jx/git-operator/resources` in a git clone of the repository before setting up the Secret

 
### Create the Git URL Secret

You need to create a `Secret` to map the git repository to the `jx-git-operator`. 

For private repositories this will also need a username and token/password to be able to clone the git repository.

```bash 
kubectl create secret generic jx-boot --from-literal=url=https://$GIT_USERNAME:$GIT_TOKEN@github.com/myowner/myrepo.git
kubectl label secret jx-boot git-operator.jenkins.io/kind=git-operator
```

You can use any name you like for the `Secret` - it will be used as the prefix for the `Job` resources that are created.

Once the secret has been created you should see in the logs of the operator pod (see below) that the git repository is cloned and a `Job` is triggered to apply the contents of git.
 
### Viewing the logs

To see the logs of the operator try:


```bash
kubectl logs -f -l app=jx-git-operator
```    

you should see it polling your git repository and triggering `Job` instances whenever a change is deteted


### Running 

You can run the `jx-git-operator` locally on the command line if you want. Actions will be created as Kubernetes Jobs even if you run the binary locally - it is just the git polling which runs locally.

Download the [x-git-operator binary](https://github.com/jenkins-x/x-git-operator/releases) for your operating system and add it to your `$PATH`.

There will be an `app` you can install soon too...
