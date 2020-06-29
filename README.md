## jx-git-operator

`jx-git-operator` is an operator which polls a git repository for changes and triggers a Kubernetes `Job` to process changes in git.

It can be used to install/upgrade any environment (development, staging, production) via some GitOps approach using some set of tools (`kubectl`, `helm`, `helmfile`, `kpt`, `kustomize` etc).

### Installing

To install the git operator using [helm 3](https://helm.sh/) then try:

Setup a namespace:

```bash 
kubectl create ns jxb
jx ns jxb
```

Then use helm to install/upgrade:
         
```bash    
helm repo add jx-labs https://storage.googleapis.com/jenkinsxio-labs-private/charts
helm install jxgo jx-labs/jx-git-operator
```

You can configure the git polling frequency via the `env.POLL_DURATION` property which supports go `time.Duration` syntax such as `10m` or `40s`

The chart defaults to using a `cluster-admin` role so it can create `Job` resources in any namespace along with any associated resources specified in a git repository at `.jx/git-operator/resources/*.yaml`

You can enable strict mode which only requires roles to read `Secret` resources in the namespace its installed and list/create `Job` resources via the `rbac.strict = true`. 

Also to avoid cluster roles use `rbac.cluster = false` which only uses a `Role` and `RoleBinding` in the install namespace.

### Setting up a repository

The git repository you wish to boot needs to have the `.jx/git-operator/job.yaml` defined to specify the Kubernetes `Job` to perform the boot job.

### Create the Git URL Secret

You need to create a `Secret` to map the git repository to the `jx-git-operator`. 

For private repositories this will also need a username and token/password to be able to clone the git repository.

```bash 
kubectl create secret generic jx-git-operator-boot --from-literal=url=https://myusername:mytoken@github.com/myowner/myrepo.git
kubectl label secret jx-git-operator-boot git-operator.jenkins.io/kind=git-operator
```

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
