module github.com/IBM/portieris/cosign-verify

go 1.16

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-containerregistry v0.5.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/opencontainers/image-spec v1.0.2-0.20190823105129-775207bd45b6 // indirect
	github.com/prometheus/common v0.18.0
	github.com/sigstore/cosign v0.2.0
	github.com/sigstore/sigstore v0.0.0-20210427115853-11e6eaab7cdc
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0

)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.4-0.20200207053602-7439e774c9e9+incompatible
	github.com/IBM/portieris/cosign-verify => ./
	github.com/googleapis/gnostic/OpenAPIv2 => github.com/googleapis/gnostic/openapiv2 v0.5.4
	github.com/sigstore/cosign => ./cosign // define path to cosign
	// google.golang.org/grpc => google.golang.org/grpc v1.29.1
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)
