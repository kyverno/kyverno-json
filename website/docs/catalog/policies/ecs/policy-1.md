
# policy-1

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/policy-1.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/policy-1.yaml
```

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
        - check:
            foo:
              /(bar)/: 10
```
