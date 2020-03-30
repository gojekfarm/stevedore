---
id: overview
title: Overview
sidebar_label: Overview
---

In a nutshell, **Stevedore** is a *wrapper over* **helm**, yet it provides powerfull features around it.

#### Values Management
Helm is the package manager for kubernetes which helps us bundle the list of yaml files as `Charts` and use them to create helm releases. We can configure the release using the values.yaml (or `helm install --values`). Using helm we are able to reduce lot of duplication in yaml files and we can reuse the bundle(Chart) with different values.yaml files.

Stevedore went one level above and started **reducing the duplication** in `values.yaml` using [overrides](overrides.md).

#### Dependencies Management
In helm, we can combine mutliple charts and [create a dependencies chart](https://helm.sh/docs/developing_charts/#managing-dependencies-with-requirements-yaml), stevedore can **build and publish the dependenices chart** on the fly based on changes in [manifest](chartspec.md) to your own hosted [chartmuseum](https://chartmuseum.com/)

#### Secrets Management
With helm, we have to either hardcode the **secrets** in `values.yaml` or store this files in some store and fetch and give to helm while installing/upgrading.

Stevedore can fetch configs on the fly from the store and do install/upgrade. It will consider values of environment (env variables) and also the process values (process variable) and give the [order of precedence](envs.md) to them. To support all kinds of store, Stevedore is following [plugin architecture](https://github.com/hashicorp/go-plugin) very similar to terraform/vault and currently we have plugin for [Consul](https://www.consul.io/), if needed please write one for your own store.

#### Release Management
With helm, we may not able to know the details of changes that are gonna get installed/upgraded as we may not know what the underlying chart has. We wanted to know the changes at `kubernetes` resource level. To achieve that we have to use tools like [helm diff](https://github.com/databus23/helm-diff#helm-diff-plugin).

Stevedore using `helm diff` plugin for the same and uses `plan and apply` style

#### Cluster Management
For helm to access cluster, we have to run tillers(for version 2) and create role, rolebinding with proper values. Stevedore can do everything in [one step](init.md)

<!--TODO
1. chartspec link is broken
2. order of precedence section
3. Can we have credits page for helm diff, hashicorp plugin
4. gif for plan and apply and attach the link here
5. init.md is broken
-->