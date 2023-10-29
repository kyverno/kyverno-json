---
tags:
- aws
- aws/s3
---
# Policy 1

## Description

Policy 1

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/aws/policy-1.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/aws/policy-1.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/aws/policy-1.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/aws/policy-1.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: Policy 1
    title.policy.kyverno.io: Policy 1
  creationTimestamp: null
  labels:
    s3.aws.tags.kyverno.io: ""
  name: test
spec:
  rules:
  - assert:
      all:
      - check:
          foo:
            /(bar)/: 10
    name: foo-bar
```
