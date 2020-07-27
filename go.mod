module github.com/jenkins-x/jx-git-operator

go 1.13

require (
	github.com/jenkins-x/jx-helpers v1.0.7
	github.com/jenkins-x/jx-kube-client v0.0.8
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.1
	github.com/stretchr/testify v1.6.0
	k8s.io/api v0.17.6
	k8s.io/apimachinery v0.17.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
)

replace (
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.16.5
)
