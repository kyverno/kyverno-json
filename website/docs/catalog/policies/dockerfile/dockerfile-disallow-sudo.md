---
tags:
- dockerfile
---
# Ensure sudo is not used in Dockerfile

## Description

This Policy ensures that sudo isn’t used.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-sudo.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-sudo.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/dockerfile/dockerfile-disallow-sudo.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-sudo.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that sudo isn’t used.
    title.policy.kyverno.io: Ensure sudo is not used in Dockerfile
  creationTimestamp: null
  labels:
    dockerfile.tags.kyverno.io: dockerfile
  name: dockerfile-disallow-sudo
spec:
  rules:
  - assert:
      all:
      - check:
          ~.(Stages[].Commands[].CmdLine[]):
            (contains(@, 'sudo')): false
        message: sudo not allowed
    name: dockerfile-disallow-sudo
```
