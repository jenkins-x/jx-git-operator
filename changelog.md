### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.155/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.155/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* just use multi-arch images (James Strachan)
* avoid including the git user/token in the logs (James Strachan)

### Chores

* upgrade go dependencies (jenkins-x-bot)
* simplify pipelines (James Strachan)
* update deps (James Strachan)
