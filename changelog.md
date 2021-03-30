### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.170/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v0.0.170/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Bug Fixes

* better support for numeric git username (James Strachan) [#158](https://github.com/jenkins-x/jx-git-operator/issues/158) 

### Chores

* avoid release CRD (James Strachan)

### Issues

* [#158](https://github.com/jenkins-x/jx-git-operator/issues/158) Allow integer type for Git username
