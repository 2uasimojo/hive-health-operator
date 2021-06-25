module github.com/openshift/hive-health-operator

go 1.16

require (
	github.com/openshift/hive/apis v0.0.0-20210624144808-697460baf215
	github.com/operator-framework/operator-sdk v0.18.2
	github.com/spf13/pflag v1.0.5
	k8s.io/api v0.20.8
	k8s.io/apimachinery v0.20.8
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.19
	k8s.io/api => k8s.io/api v0.20.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.8
	k8s.io/client-go => k8s.io/client-go v0.20.8
)
