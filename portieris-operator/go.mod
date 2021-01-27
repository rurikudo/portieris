module github.com/rurikudo/portieris/portieris-operator

go 1.13

require (
	github.com/IBM/portieris v0.10.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/openshift/api v3.9.0+incompatible
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.3
)

replace github.com/rurikudo/portieris/portieris-operator => ./
