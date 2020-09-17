module github.com/jenkins-x/jx-git-operator

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.9.8 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jenkins-x/jx-helpers v1.0.45
	github.com/jenkins-x/jx-kube-client v0.0.8
	github.com/jenkins-x/jx-logging v0.0.11
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/text v0.3.3 // indirect
	k8s.io/api v0.17.11
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.2
	k8s.io/client-go => k8s.io/client-go v0.16.5
)
