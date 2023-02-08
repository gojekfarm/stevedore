module github.com/gojek/stevedore

go 1.14

replace github.com/micro/go-micro => github.com/micro/go-micro v1.6.0

require (
	github.com/Masterminds/semver v1.5.0
	github.com/aryann/difflib v0.0.0-20210328193216-ff5ff6dc229b // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/chartmuseum/helm-push v0.7.1
	github.com/cucumber/godog v0.12.1
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/databus23/helm-diff v3.1.1+incompatible
	github.com/fatih/color v1.12.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.9
	github.com/hashicorp/consul/api v1.3.0 // indirect
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-plugin v1.0.1-0.20190430211030-5692942914bb
	github.com/imdario/mergo v0.3.12
	github.com/json-iterator/go v1.1.12
	github.com/manifoldco/promptui v0.8.0
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/micro/go-micro v1.6.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	golang.org/x/term v0.4.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.0
	gopkg.in/h2non/gock.v1 v1.1.2
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	helm.sh/helm/v3 v3.11.1
	k8s.io/cli-runtime v0.26.0
	k8s.io/client-go v0.26.0
	k8s.io/helm v2.17.0+incompatible // indirect
)
