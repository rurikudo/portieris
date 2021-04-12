module github.com/rurikudo/portieris

go 1.16

require (
	github.com/IBM/go-sdk-core v0.9.0
	github.com/IBM/go-sdk-core/v4 v4.10.0
	github.com/Shopify/logrus-bugsnag v0.0.0-20171204204709-577dee27f20d // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/bitly/go-hostpool v0.1.0 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bugsnag/bugsnag-go v1.5.3 // indirect
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/cloudflare/cfssl v1.5.0 // indirect
	github.com/containers/image/v5 v5.9.0
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/go v1.5.1-1 // indirect
	github.com/form3tech-oss/jwt-go v3.2.1+incompatible // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-containerregistry v0.4.1
	github.com/googleapis/gnostic v0.5.4
	github.com/gorilla/mux v1.7.4
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jinzhu/gorm v1.9.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kubernetes/apiextensions-apiserver v0.0.0-20181121072900-e8a638592964
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/prometheus/client_golang v1.10.0
	github.com/prometheus/common v0.18.0
	github.com/satori/go.uuid v1.2.0
	github.com/sigstore/cosign v0.2.0
	github.com/sigstore/sigstore v0.0.0-20210405172749-e614ea31ba83
	github.com/stretchr/testify v1.7.0
	github.com/theupdateframework/notary v0.6.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	google.golang.org/grpc/examples v0.0.0-20210408231144-1d1bbb55b381 // indirect
	gopkg.in/dancannon/gorethink.v3 v3.0.5 // indirect
	gopkg.in/fatih/pool.v2 v2.0.0 // indirect
	gopkg.in/gorethink/gorethink.v3 v3.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery v0.19.0
	k8s.io/client-go v0.19.0
)

replace (
	github.com/IBM/portieris => ./
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.4-0.20200207053602-7439e774c9e9+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
	github.com/googleapis/gnostic/OpenAPIv2 =>  github.com/googleapis/gnostic/openapiv2 v0.5.4
	k8s.io/api => k8s.io/api v0.19.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.0
	k8s.io/client-go => k8s.io/client-go v0.19.0
)