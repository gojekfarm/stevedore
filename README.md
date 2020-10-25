# [Stevedore](https://en.wikipedia.org/wiki/Stevedore)

[![Test](https://github.com/gojekfarm/stevedore/workflows/Check/badge.svg)](https://github.com/gojekfarm/stevedore/actions?query=workflow%3ACheck+branch%3Amaster)

![logo](/logo/logo_readme.png)

Stevedore a tool to load the cluster with containers for kubernetes to orchestrate. It is a wrapper on Helm, gets all the features of helm, defines a workflow through which helm charts can be deployed and managed. It offers,

- **Declarative way** to mange dependencies.
- Terraform style **plan and apply** to see what exactly going to change
- Ability to specify overrides for any given environment at ease
- Support storing of configurations multiple config store (Plugin support for custom stores)

**TLDR;** Kubernetes follows a declarative style of configuring infrastructure. Stevedore provides a way through which helm interactions can be done declarative through a YAML configuration.

![demo](demo.gif)

## Table of Contents

1. [Getting Started](#getting-started)

   1.1 [Stevedore Context](#stevedore-context)

   1.2 [Manifest](#Manifest)

   1.3 [Plan](#plan)

   1.4 [Apply](#apply)

   1.5 [Using Override](#using-override)

   1.6 [Using Env](#using-env)

2. [Terminology](#terminology)

3. [Development](#development)

## Getting Started

Example guide for installing 'redis' helm chart via the stevedore into minikube.

### Stevedore Context

First, add the kubernetes cluster to stevedore context, in our case its minikube

```bash
$ stevedore config add-context \
    --name minikube \
    --environment local \
    --environment-type dev \
    --kube-context minikube \
    --type local

Successfully added the below context:
Name                    minikube
Type                    local
Environment             local
Environment Type        dev
Kubernetes Context      minikube

Successfully switched to context: minikube
```

### Manifest

create a stevedore manifest file and save it as `redis.yaml`

```yaml
kind: StevedoreManifest
version: 2
deployTo:
  - contextName: minikube
    environmentType: dev
spec:
  - release:
      name: redis
      namespace: default
      chart: stable/redis
      values:
        password: password
```

### Plan

```bash
$ stevedore plan -f redis.yaml
```

stevedore will show the detailed plan on what are the resources which will be created

```bash
## partial output of plan command
Release changes:
RELEASE         MANIFEST CHANGES
redis           +ConfigMap/redis
(redis.yaml)    +ConfigMap/redis-health
                +Secret/redis
                +Service/redis-headless
                +Service/redis-master
                +Service/redis-slave
                +StatefulSet/redis-master
                +StatefulSet/redis-slave

File changes:
FILENAME        ADDITIONS       MODIFICATIONS   DESTRUCTIONS
redis.yaml      8               0               0
```

> plan **will not** modify the actual resources.
>
> to use the output of the plan to apply command, specify the artifact folder via additional artifact path --artifact/-a
> stevedore plan -f redis.yaml -a out

### Apply

to, install / update the changes

```bash
$ stevedore apply -f redis.yaml
```

to install the / update the changes planned via plan command specify the files from the artifact directory

```bash
$ stevedore apply -f out/redis.yaml
```

to proceed further and persist the changes, confirm the action

```bash
Context Details:
------------------
Name: minikube
Type: local
Environment: local
Kubernetes Context: minikube
Environment Type: dev
KubeConfig File:
------------------
Confirm to apply: [y/N]
```

to auto apply specify --yes

### Using override

Stevedore offers easy way to manage overrides for different environment.

consider the following requirements for four different environments (dev, test, integration and production) where redis
will be installed.

1. dev, test environment may not need persistent to be enabled.

```yaml
kind: StevedoreOverride
version: 2
spec:
  - matches:
    environmentType: dev
    values:
      master:
        persistence:
          enabled: false
      slave:
        persistence:
          enabled: false
  - matches:
    environmentType: test
    values:
      master:
        persistence:
          enabled: false
      slave:
        persistence:
          enabled: false
```

2. integration environment resource request will be little higher than dev and test environment

```yaml
kind: StevedoreOverride
version: 2
spec:
  - matches:
    environmentType: integration
    values:
      master:
        resources:
          requests:
            memory: 256Mi
            cpu: 100m
      slave:
        resources:
          requests:
            memory: 256Mi
            cpu: 100m
```

3. production environment needs persistent to be enabled and has the highest value for resource quota

```yaml
kind: StevedoreOverride
version: 2
spec:
  - matches:
    environmentType: production
    values:
      master:
        resources:
          requests:
            memory: 1024Mi
            cpu: 1000m
        persistence:
          enabled: false
      slave:
        resources:
          requests:
            memory: 1024Mi
            cpu: 1000m
        persistence:
          enabled: false
```

to use the override feature, save these overrides as `override.yaml` and during plan and apply pass on the overrides

```bash
$ stevedore plan -f redis.yaml -o override.yaml
$ stevedore apply -f redis.yaml -o override.yaml
```

### Using Env

Similar to overrides, stevedore provides an easy way to manage environment, cluster specific overrides

to provide a different password for different environment templatize the password field as `${PASSWORD}`

```yaml
kind: StevedoreManifest
version: 2
deployTo:
  - contextName: minikube
    environmentType: dev
spec:
  - release:
      name: redis
      namespace: default
      chart: stable/redis
      values:
        password: ${REDIS_PASSWORD}
```

now, you can define an environment file or export an environment variable for the stevedore to pick it up.

```yaml
kind: StevedoreEnv
version: 2
spec:
  - matches:
      environmentType: dev
    env:
      REDIS_PASSWORD: dev
---
kind: StevedoreEnv
version: 2
spec:
  - matches:
      environmentType: test
    env:
      REDIS_PASSWORD: test
```

## Terminology

**StevedoreManifest** use this to define the release manifest which is interpreted by the stevedore and perform install / upgrade

**StevedoreEnv** use this to define environment/context/cluster specific value (eg., database connection string, password, replica count etc.,)

**StevedoreOverride** use this to define overrides for an environment/context/cluster

## Development

### Requirements

1. `go` version: `1.14`
2. Clone the repo to `$GOPATH/src/github.com/gojekfarm/stevedore`
3. Run `make build` to build dependencies, format, vet, lint, test and compile
4. Run `make compile` to compile for local development machine (use `compile-linux` to compile for linux os amd64 architecture)
5. Run `make install` to install stevedore \$GOPATH/bin

## Credits

Logo designed by [Kartik Narayanan](https://varnaturika.art) checkout his various work [here](https://varnaturika.art),
you can reach out him at [@hajaarfunda](https://twitter.com/hajaarfunda)

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
