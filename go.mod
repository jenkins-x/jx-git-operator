module github.com/jenkins-x/jx-git-operator

go 1.15

require (
	github.com/fatih/color v1.10.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11
	github.com/jenkins-x/jx-helpers/v3 v3.0.88
	github.com/jenkins-x/jx-kube-client/v3 v3.0.2
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.3.2
	github.com/sirupsen/logrus v1.7.1 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20210113233702-8566a335510f // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)
