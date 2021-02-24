### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.161/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.161/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* lets avoid adding user/pwd to commit url annotation (James Strachan)
* lets annotate the job with commit information (James Strachan)
* disable the use of custom file names (James Strachan)

### Chores

* upgrade deps (James Strachan)
* fix failing test (James Strachan)
* update deps (James Strachan)
