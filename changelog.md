### Linux

```shell
curl -L https://github.com/jenkins-x/jx-git-operator/releases/download/v{{ .Chart.Version }}/jx-git-operator-linux-amd64.tar.gz | tar xzv 
sudo mv jx-git-operator /usr/local/bin
```

### macOS

```shell
curl -L  https://github.com/jenkins-x/jx-git-operator/releases/download/v{{ .Chart.Version }}/jx-git-operator-darwin-amd64.tar.gz | tar xzv
sudo mv jx-git-operator /usr/local/bin
```

## Changes

### Chores

* fix gha (James Rawlings)
