---
id: matches
title: Matches
sidebar_label: Matches
---

For creating application/release values, stevedore will look for the matches based on the condition.

Each matches can have one or more conditions. Below are types of conditions with its weightage:

```yaml
environmentType (1)
environment (2)
contextType (4)
contextName (8)
applicationName (16)
```

Weightage for a matches is the sum of each condition's weightage.

If there are multiple matches, stevedore will merge based on the weightage.

Consider the below mainfest:

```yaml
deployTo:
  - test-components
spec:
    release:
      name: some-service
      namespace: default
      chart: stable/postgres
      values:
        db: {}
    ...
```

and below overrides,

```yaml
kind: StevedoreOverride
version: '2'
spec:
- matches:
    contextType: components
  values:
    db:
      persistence:
        size: 10Gi
      user:
        name: test
- matches:
    applicationName: some-service
  values:
    db:
      persistence:
        size: 5Gi
```

the final value in the manifest will be

```yaml
deployTo:
  - cluster-1
spec:
  release:
    name: some-service
    namespace: default
    chart: stable/postgres
  values:
    db:
      persistence:
        size: 5Gi
      user:
        name: test
```