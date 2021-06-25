module github.com/openshift/hive-health-operator

go 1.16

require (
	github.com/Azure/go-autorest/autorest v0.11.12 // indirect
	github.com/openshift/hive/apis v0.0.0-20210624144808-697460baf215
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.8.2
	github.com/Azure/go-autorest/autorest/azure/cli => github.com/Azure/go-autorest/autorest/azure/cli v0.3.1
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/Azure/go-autorest/autorest/validation => github.com/Azure/go-autorest/autorest/validation v0.2.1-0.20191028180845-3492b2aff503
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.1
	k8s.io/api => k8s.io/api v0.20.8
	k8s.io/client-go => k8s.io/client-go v0.20.0
)
