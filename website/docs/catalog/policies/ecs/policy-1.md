
# policy-1

## Description

None

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/catalog/ecs/policy-1.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: Policy
metadata:
  creationTimestamp: null
  name: test
spec:
  rules:
  - name: foo-bar
    validate:
      assert:
        all:
        - foo:
            /(bar)/: 10
```