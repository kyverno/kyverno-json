---
tags:
- dockerfile
---
# Dockerfile expose port 22 not allowed

## Description

This Policy ensures that port 22 is not exposed in Dockerfile.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-deny-expose-22.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-deny-expose-22.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/dockerfile/dockerfile-deny-expose-22.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-deny-expose-22.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that port 22 is not exposed
      in Dockerfile.
    title.policy.kyverno.io: Dockerfile expose port 22 not allowed
  creationTimestamp: null
  labels:
    dockerfile.tags.kyverno.io: dockerfile
  name: dockerfile-deny-expose-port-22
spec:
  rules:
  - assert:
      all:
      - check:
          ~.(Stages[].Commands[?Name=='EXPOSE'][]):
            (contains(Ports, '22') || contains(Ports, '22/TCP')): false
        message: Port 22 exposure is not allowed
    name: check-port-exposure
```
