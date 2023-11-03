---
tags:
- dockerfile
---
# Ensure apt is not used in Dockerfile

## Description

This Policy ensures that apt isnt used but apt-get can be used as apt interface is less stable than apt-get and so this preferred.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-apt.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-apt.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/dockerfile/dockerfile-disallow-apt.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-disallow-apt.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that apt isnt used but apt-get
      can be used as apt interface is less stable than apt-get and so this preferred.
    title.policy.kyverno.io: Ensure apt is not used in Dockerfile
  creationTimestamp: null
  labels:
    dockerfile.tags.kyverno.io: dockerfile
  name: dockerfile-disallow-apt
spec:
  rules:
  - assert:
      any:
      - check:
          ~.(Stages[].Commands[].CmdLine[]):
            (contains(@, 'apt ')): false
        message: apt not allowed
    name: dockerfile-disallow-apt
```
