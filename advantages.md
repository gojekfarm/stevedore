### Advantages of Helm
- Official Package manager for Kubernetes
- No need to know kubernetes at all (just values.yaml and helm commands)
- Easy to upgrade/downgrade releases in one step
- Easy to share the helm charts
- Chart versioning
- Charts based on other dependencies charts
- Can write gotemplate functions,if-else,range,etc 
- Helm provides a hook mechanism like pre-install,pre-delete,pre-upgrade,post-rollback,etc by which we can have more customization in helm charts

### Advantages of Stevedore
- All advantages of Helm (as stevedore is wrapper over helm)
- Declarative way
- Gitops
- No duplication of values.yaml for multiple envs/services
- Plan and Apply approach. We know what exactly is going to get applied
- Configuration can be from env variables, consul kv store. (as it is plugin architecture, we can other storages like etcd, vault)
- Battle tested in gojek to handle 100s of helm releases across multiple envs and multiple countries
- Merge yaml files based on weights and produce final values.yaml for helm release
- Supports both helm2 and helm3
