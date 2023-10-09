# Match and exclude

Both [Kyverno policies](https://kyverno.io/docs/kyverno-policies/) and `kyverno-json` policies can match and exclude *resources* when being evaluated.

[Kyverno policies](https://kyverno.io/docs/kyverno-policies/) use [Kubernetes](https://kubernetes.io) specific constructs for that matter that didn't map well with arbitrary payloads.

`kyverno-json` uses [assertion trees](./assertion-trees.md) to implement `match` and `exclude` statements:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: Policy
metadata:
  name: required-s3-tags
spec:
  rules:
    - name: require-team-tag
      match:
        any:
        - type: aws_s3_bucket
      exclude:
        any:
        - name: bypass-me
      validate:
        assert:
          all:
          - values:
              tags:
                Team: ?*
```

In the example above, every *resource* having `type: aws_s3_bucket` will match, and *resources* having `name: bypass-me` will be excluded.
