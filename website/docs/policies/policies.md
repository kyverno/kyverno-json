# Policy Structure

Kyverno policies are [Kubernetes resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and can be easily managed via Kubernetes APIs, GitOps workflows, and other existing tools.

Policies that apply to JSON payload have a few differences from Kyverno policies that are applied to Kubernetes resources at admission controls.

## Resource Scope

Policies that apply to JSON payloads are always cluster-wide resources.

## API Group and Kind

`kyverno-json` policies belong to the `json.kyverno.io` group and can only be of kind `ValidatingPolicy`.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
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

## Policy Rules

A policy can have multiple rules, and rules are processed in order. Evaluation stops at the first rule that fails.

## Match and Exclude

Policies that apply to JSON payloads use [assertion trees](./asserts.md) in both the `match`/`exclude` declarations as well as the `validate` rule declaration.

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: required-s3-tags
spec:
  rules:
    - name: require-team-tag
      identifier: address
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

In the example above, every *resource* having `type: aws_s3_bucket` will match, and *payloads* having `name: bypass-me` will be excluded.

## Identifying Payload Entries

A policy rule can contain an optional `identifier` which declares the path to the payload element that uniquely identifies each entry.

## Context Entries

A policy rule can contain optional `context` entries that are made available to the rule via bindings:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: required-s3-tags
spec:
  rules:
    - name: require-team-tag
      match:
        any:
        - type: aws_s3_bucket
      context:
      # creates a `expectedTeam` binding automatically
      - name: expectedTeam
        variable: Kyverno
      validate:
        message: Bucket `{{ name }}` does not have the required Team tag {{ $expectedTeam }}
        assert:
          all:
          - values:
              tags:
                # use the `$expectedTeam` binding coming from the context
                Team: ($expectedTeam)
```

## No `forEach`, `pattern operators`, `anchors`, or `wildcards`

The use of [assertion trees](./asserts.md) addresses some features of Kyverno policies that apply to Kubernetes resources.

Specifically, [forEach](https://kyverno.io/docs/writing-policies/validate/#foreach), [pattern operators](https://kyverno.io/docs/writing-policies/validate/#operators), [anchors](https://kyverno.io/docs/writing-policies/validate/#anchors), or [wildcards](https://kyverno.io/docs/writing-policies/validate/#wildcards) are not supported for policies that apply to JSON resources. Instead, [assertion trees](./asserts.md) with [JMESPath](../jp.md) expressions are used to achieve the same powerful features.