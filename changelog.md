### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.168/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.168/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* add quotes around the OCI secret values (James Strachan)
* add badges (James Strachan)
