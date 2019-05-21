module github.com/gowait

require (
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20180628040859-072894a440bd // indirect
	k8s.io/apimachinery v0.0.0-20180621070125-103fd098999d // indirect
	k8s.io/client-go v8.0.0+incompatible // indirect
)

replace config => ./config
