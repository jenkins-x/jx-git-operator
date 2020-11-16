module github.com/jenkins-x/jx-git-operator

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.9.8 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.3 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.3 // indirect
	github.com/google/uuid v1.1.2
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jenkins-x/jx-api/v3 v3.0.3 // indirect
	github.com/jenkins-x/jx-helpers/v3 v3.0.18
	github.com/jenkins-x/jx-kube-client/v3 v3.0.1
	github.com/jenkins-x/jx-logging/v3 v3.0.2
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.1.2
	github.com/stretchr/testify v1.6.1
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v0.19.4
	sigs.k8s.io/structured-merge-diff/v4 v4.0.2 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2
