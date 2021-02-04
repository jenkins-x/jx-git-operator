### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.154/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.154/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* just use multi-arch images (James Strachan)
* avoid including the git user/token in the logs (James Strachan)

### Chores

* simplify pipelines (James Strachan)
* update deps (James Strachan)
* deps: bump https://github.com/jenkins-x/jx-cli-base-image.git to 0.0.43 (jenkins-x-bot)
