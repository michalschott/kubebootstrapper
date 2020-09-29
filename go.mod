module github.com/michalschott/kubebootstrapper

go 1.15

replace (
	github.com/michalschott/kubebootstrapper/cli/cmd => ./cli/cmd
	github.com/michalschott/kubebootstrapper/pkg/ptr => ./pkg/ptr
	github.com/michalschott/kubebootstrapper/pkg/version => ./pkg/version
)

require (
	github.com/briandowns/spinner v1.11.1
	github.com/fatih/color v1.9.0
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.0.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
)
