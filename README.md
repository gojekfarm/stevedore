# [Stevedore](https://en.wikipedia.org/wiki/Stevedore)

Stevedore a tool to load the cluster with containers for kubernetes to orchestrate. It is a wrapper on Helm, gets all the features of helm, defines a workflow through which helm charts can be deployed and managed. It offers,

- **Declartive way** to mange dependencies.
- Terraform style **plan and apply** to see what exactly going to change
- Ability to specify overrides for any given environment at ease
- Support storing of configurations multiple config store (Plugin support for custom stores)

**TLDR;** Kubernetes follows a declarative style of configuring infrastructure. Stevedore provides a way through which helm interactions can be done declaratively through a YAML configuration.

## Table of Contents

1. [Command Usage](#commands)

   - [config](#config)
   - [render](#render)
   - [plan](#plan)
   - [apply](#apply)

2. [Example](#example)

   - [Manifest](#manifest)
   - [Override](#override)
   - [Env](#env)

3. [Usage](#usage)

   - [Manifest](#manifest)
   - [Override](#override)
   - [Envs](#envs)

4. [Development](#development)

## Commands

### config

```
Manage stevedore config

Usage:
  stevedore config [command]

Available Commands:
  add-context    Adds context to stevedore config file
  delete-context Delete the specified context from the stevedore config
  get-contexts   Describe one or many contexts
  rename-context Renames a context from the stevedore config file
  show-context   Shows the current-context
  use-context    Sets the current-context in a stevedore config file
  view           Display complete configuration

Flags:
  -h, --help   help for config

Global Flags:
      --config string       config file (default "${HOME}/.stevedore/config")
      --kubeconfig string   path to kubeconfig file (default: ${HOME}/.kube/config)
      --log-level string    set the logger level (default "error")

Use "stevedore config [command] --help" for more information about a command.
```

### render

```
Validate and render stevedore yaml(s)

Usage:
  stevedore render [flags]

Flags:
  -a, --artifact-path string    Stevedore artifact path (folder) to save the output as artifact
  -e, --env-path string         Stevedore env path (can be yaml file or folder)
  -h, --help                    help for render
  -f, --manifest-path string    Stevedore manifest path (can be yaml file or folder)
  -o, --overrides-path string   Stevedore overrides path (can be yaml file or folder)

Global Flags:
      --config string       config file (default "${HOME}/.stevedore/config")
      --kubeconfig string   path to kubeconfig file (default: ${HOME}/.kube/config)
```

      --log-level string    set the logger level (default "error")

### plan

```
Validate and plan stevedore yaml(s)

Usage:
  stevedore plan [flags]

Flags:
  -a, --artifact-path string    Stevedore artifact path (folder) to save the output as artifact
  -e, --env-path string         Stevedore env path (can be yaml file or folder)
  -h, --help                    help for plan
  -f, --manifest-path string    Stevedore manifest path (can be yaml file or folder)
  -o, --overrides-path string   Stevedore overrides path (can be yaml file or folder)

Global Flags:
      --config string       config file (default "${HOME}/.stevedore/config")
      --kubeconfig string   path to kubeconfig file (default: ${HOME}/.kube/config)
      --log-level string    set the logger level (default "error")
```

### apply

```
Validate and apply stevedore yaml(s)

Usage:
  stevedore apply [flags]

Flags:
  -a, --artifact-path string    Stevedore artifact path (folder) to save the output as artifact
  -e, --env-path string         Stevedore env path (can be yaml file or folder)
  -h, --help                    help for apply
  -f, --manifest-path string    Stevedore manifest path (can be yaml file or folder)
  -o, --overrides-path string   Stevedore overrides path (can be yaml file or folder)

Global Flags:
      --config string       config file (default "${HOME}/.stevedore/config")
      --kubeconfig string   path to kubeconfig file (default: ${HOME}/.kube/config)
      --log-level string    set the logger level (default "error")
```

## Development

**Requirements**

1. `go` version: `1.13`
3. Clone the repo to `$GOPATH/src/github.com/gojekfarm/stevedore`
4. Run `make build` to build dependencies, format, vet, lint, test and compile
5. Run `make compile` to compile for local development machine (use `compile-linux` to compile for linux os amd64 architecture)
6. Run `make install` to install stevedore $GOPATH/bin

## License

Copyright 2018-2020, GO-JEK Tech (http://gojek.tech)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
