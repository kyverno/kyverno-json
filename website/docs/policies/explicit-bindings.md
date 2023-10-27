# Explicit bindings

Sometimes it can be useful to refer to a parent node in the assertion tree.

This is possible to add an explicit binding at every node in the tree by appending the `@binding_name` to the key.

Given the input document:

```yaml
foo:
  bar: 4
  bat: 6
```

The following policy will compute a sum and bind the result to the `sum` binding. A descendant can then use `$sum` and use it:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: Policy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar
      validate:
        assert:
          all:
          - foo:
              # evaluate expression `(bar + bat)` and bind it to `sum`
              (bar + bat)@sum:
                # get the `$sum` binding and compare it against `10`
                ($sum): 10
```

All binding are available to descendants, if a descendant creates a binding with a name that already exists the binding will be overriden for descendants only and it doesn't affect the bindings at upper levels in the tree.

In other words, a node in the tree always sees bindings that are definied in the parents and if a name is reused, the first binding with the given name wins when winding up the tree.

As a consequence, the policy below is perfectly valid:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: Policy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar
      validate:
        assert:
          all:
          - foo:
              (bar + bat)@sum:
                ($sum + $sum)@sum:
                  ($sum): 20
                ($sum): 10
```

Note that all context entries are made available to the rule via bindings:

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
      context:
      # creates a `expectedTeam` binding automatically
      - name: expectedTeam
        variable: Kyverno
      validate:
        message: Bucket `{{ name }}` ({{ address }}) does not have the required Team tag {{ $expectedTeam }}
        assert:
          all:
          - values:
              tags:
                # use the `$expectedTeam` binding coming from the context
                Team: ($expectedTeam)
```

Finally, we can always access the current payload, policy and rule being evaluated using the builtin `$payload`, `$policy` and `$rule` bindings. No protection is made to prevent you from overriding those bindings though.
