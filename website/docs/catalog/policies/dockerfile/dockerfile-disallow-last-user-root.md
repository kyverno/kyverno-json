---
tags:
- dockerfile
---
# Dockerfile last user is not allowed to be root

## Description

This Policy ensures that last user in Dockerfile is not root.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-last-user-root.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-last-user-root.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/dockerfile/dockerfile-disallow-last-user-root.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-last-user-root.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that last user in Dockerfile
      is not root.
    title.policy.kyverno.io: Dockerfile last user is not allowed to be root
  creationTimestamp: null
  labels:
    dockerfile.tags.kyverno.io: dockerfile
  name: dockerfile-disallow-last-user-root
spec:
  rules:
  - assert:
      all:
      - check:
          ((Stages[].Commands[?Name == 'USER'][])[-1].User == 'root'): false
        message: Last user root not allowed
    name: check-disallow-last-user-root
```
