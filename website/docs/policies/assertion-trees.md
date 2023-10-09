# Assertion trees

[Kyverno policies](https://kyverno.io/docs/kyverno-policies/) started with a declarative approach but slowly adopted the imperative approach too, because of the limitations in the implemented declarative approach.

`kyverno-json` tries to be as declarative as possible, for now `forEach`, pattern operators, anchors and wildcards are not supported are not supported.
Hopefully we won't need to adopt an imperative approach anymore.

Assertion trees can be used to express complex and dynamic conditions by using [jmespath](https://jmespath.site) expressions.

Those expressions represent projections of the being analysed *resource* and the result of this projection is passed to descendants for further analysis.

All comparisons happen in the leaves of the assertion tree.

**Example**:

Given the input payload below:

```yaml
foo:
  baz: true
  bar: 4
  bat: 6
```

It is possible to write a validation rule like this:

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
          -
            # project field `foo` onto itself, the content of `foo` becomes the current object for descendants
            foo:

              # evaluate expression `(bar > `3`)`, the result becomes the current object for descendants (in this case the result will be a simple boolean)
              # then we hit the `true` leaf, comparison happens and we expect the current value to be `true`
              (bar > `3`): true

              # evaluate expression `(!baz)`, the result becomes the current object for descendants (in this case the result will be a simple boolean)
              # then we hit the `true` leaf, comparison happens and we expect the current value to be `false`
              (!baz): false

              # evaluate expression `(bar + bat)`, the result becomes the current object for descendants (in this case the result will be a number)
              # then we hit the `10` leaf, comparison happens and we expect the current value to be `10`
              (bar + bat): 10
```
