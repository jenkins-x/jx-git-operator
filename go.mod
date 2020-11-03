module github.com/jenkins-x/jx-git-operator

go 1.13

require (
	github.com/jenkins-x/jx-helpers/v3 v3.0.14
	github.com/jenkins-x/jx-kube-client/v3 v3.0.1
	github.com/jenkins-x/jx-logging/v3 v3.0.2
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.2
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.2
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
