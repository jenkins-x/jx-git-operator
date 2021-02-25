### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.163/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.163/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* bring back support for custom boot file names (James Strachan)
* use correct git commit URL (James Strachan)
