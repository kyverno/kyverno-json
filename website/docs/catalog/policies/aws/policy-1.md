---
tags:
- aws
- aws/s3
---
# policy-1

## Description

None

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/catalog/aws/policy-1.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: Policy
metadata:
  annotations:
    description.catalog.kyverno.io: Policy 1
    title.catalog.kyverno.io: Policy 1
  creationTimestamp: null
  labels:
    s3.aws.tags.kyverno.io: ""
  name: test
spec:
  rules:
  - name: foo-bar
    validate:
      assert:
        all:
        - check:
            foo:
              /(bar)/: 10
```