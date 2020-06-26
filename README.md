## jx-git-operator

`jx-git-operator` is an operator which polls a git repository for changes and triggers a Kubernetes `Job` to process changes in git.

It can be used to install/upgrade any environment (development, staging, production) via some GitOps approach using some set of tools (`kubectl`, `helm`, `helmfile`, `kpt`, `kustomize` etc).

### Setting up a repository

### Create the Git URL Secret

You need to create a `Secret` to map the git reopsitory to the `jx-git-operator`. 

For private repositories this will also need a username and token/password to be able to clone the git repository.

```bash 
kubectl create secret generic jx-git-operator-boot --from-literal=url=https://myusername:mytoken@github.com/myowner/myrepo.git
kubectl label secret jx-git-operator-boot git-operator.jenkins.io/kind=git-operator
```
