module github.com/rurikudo/portieris/portieris-operator

go 1.13

require (
	github.com/IBM/portieris v0.10.0
	github.com/go-logr/logr v0.2.1-0.20200730175230-ee2de8da5be6
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/openshift/api v3.9.0+incompatible
	k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/kubectl v0.19.0 // indirect
	sigs.k8s.io/controller-runtime v0.6.3
)

replace github.com/rurikudo/portieris/portieris-operator => ./

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.4-0.20200207053602-7439e774c9e9+incompatible
	github.com/go-logr/logr => github.com/go-logr/logr v0.1.0
	k8s.io/api => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.6
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.6
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.io/kubectl => k8s.io/kubectl v0.18.6

)

replace github.com/openshift/api => github.com/openshift/api v0.0.0-20190924102528-32369d4db2ad // Required until https://github.com/operator-framework/operator-lifecycle-manager/pull/1241 is resolved
