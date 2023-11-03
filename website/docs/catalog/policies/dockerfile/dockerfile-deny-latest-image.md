---
tags:
- dockerfile
---
# Dockerfile latest image tag not allowed

## Description

This Policy ensures that no image uses the latest tag in Dockerfile.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-deny-latest-image.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-deny-latest-image.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/dockerfile/dockerfile-deny-latest-image.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/dockerfile/dockerfile-deny-latest-image.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that no image uses the latest
      tag in Dockerfile.
    title.policy.kyverno.io: Dockerfile latest image tag not allowed
  creationTimestamp: null
  labels:
    dockerfile.tags.kyverno.io: dockerfile
  name: dockerfile-deny-latest-image-tag
spec:
  rules:
  - assert:
      all:
      - check:
          ~.(Stages[].From.Image):
            (contains(@, ':latest')): false
        message: Latest tag is not allowed
    name: check-latest-tag
```
