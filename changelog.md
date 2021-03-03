### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.164/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.164/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* disable the custom files unless an env var is specified (James Strachan)
* bring back support for custom boot file names (James Strachan)

### Chores

* upgrade deps (James Strachan)
