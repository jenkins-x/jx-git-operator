module github.com/jenkins-x/jx-git-operator

go 1.15

require (
	cloud.google.com/go v0.76.0 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.11 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/gophercloud/gophercloud v0.15.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jenkins-x/jx-api/v4 v4.0.24 // indirect
	github.com/jenkins-x/jx-helpers/v3 v3.0.74
	github.com/jenkins-x/jx-kube-client/v3 v3.0.2
	github.com/jenkins-x/jx-logging/v3 v3.0.3
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-envconfig v0.3.2 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/oauth2 v0.0.0-20210201163806-010130855d6c // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.5.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210113233702-8566a335510f // indirect
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
	k8s.io/client-go => k8s.io/client-go v0.20.2
)
