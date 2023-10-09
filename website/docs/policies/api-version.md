# Api version and kind

Both [Kyverno policies](https://kyverno.io/docs/kyverno-policies/) and `kyverno-json` policies are defined using [Kubernetes](https://kubernetes.io) manifests.

They don't use the same `apiVersion` and `kind` though.

[Kyverno policies](https://kyverno.io/docs/kyverno-policies/) belong to the `kyverno.io` group, exist in multiple versions (`v1`, `v2beta1`) and can be of kind `Policy` or `ClusterPolicy`.

`kyverno-json` policies belong to the `json.kyverno.io` group, exist only in `v1alpha1` version and can only be of kind `Policy`.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: Policy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar-4
      validate:
        assert:
          all:
          - foo:
              bar: 4
```

The concept of clustered vs namespaced resources exist only in the [Kubernetes](https://kubernetes.io) world and it didn't make sense to reproduce the same pattern in `kyverno-json`.
