# Assertion trees

Assertion trees can be used to apply complex and dynamic conditional checks using [JMESPath](../jp.md) expressions.

## Assert 

An `assert` declaration contains an `any` or `all` list in which each entry contains a:

* `check`: the assertion check
* `message`: an optional message

A check can contain one or more JMESPath expressions. Expressions represent projections of selected data in the JSON *payload* and the result of this projection is passed to descendants for further analysis.

All comparisons happen in the leaves of the assertion tree.

**A simple example**:

This policy checks that a pod does not use the default service account:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: assert-sample
spec:
  rules:
    - name: foo-bar
      match:
        all:
        - apiVersion: v1
          kind: Pod
      assert:
        all:
        - message: "serviceAccountName 'default' is not allowed"
          check:
            spec:
              (serviceAccountName == 'default'): false
```

**A detailed example**:

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
kind: ValidatingPolicy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar-4
      validate:
        assert:
          all:
          - message: "..."
            check:
              # project field `foo` onto itself, the content of `foo` becomes the current object for descendants
              foo:

                # evaluate expression `(bar > `3`)`, the boolean result becomes the current object for descendants
                # the `true` leaf is compared with the current value `true`
                (bar > `3`): true

                # evaluate expression `(!baz)`, the boolean result becomes the current object for descendants
                # the leaf `false` is compared with the current value `false`
                (!baz): false

                # evaluate expression `(bar + bat)`, the numeric result becomes the current object for descendants
                # the leaf `10` is compared with the current value `10`
                (bar + bat): 10
```

## Iterating with Projection Modifiers

Assertion tree expressions support modifiers to influence the way projected values are processed.

The `~` modifier applies to arrays and maps, it mean the input array or map elements will be processed individually by descendants.

When the `~` modifier is not used, descendants receive the whole array, not each individual element.

Consider the following input document:

```yaml
foo:
  bar:
  - 1
  - 2
  - 3
```

The policy below does not use the `~` modifier and `foo.bar` array is compared against the expected array:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar
      validate:
        assert:
          all:
          - foo:
              # the content of the `bar` field will be compared against `[1, 2, 3]`
              bar:
              - 1
              - 2
              - 3
```

With the `~` modifier, we can apply descendant assertions to all elements in the array individually.
The policy below ensures that all elements in the input array are `< 5`:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar
      validate:
        assert:
          all:
          - foo:
              # with the `~` modifier all elements in the `[1, 2, 3]` array are processed individually and passed to descendants
              ~.bar:
                # the expression `(@ < `5`)` is evaluated for every element and the result is expected to be `true`
                (@ < `5`): true
```

The `~` modifier supports binding the index of the element being processed to a named binding with the following syntax `~index_name.bar`. When this is used, we can access the element index in descendants with `$index_name`.

When used with a map, the named binding receives the key of the element being processed.

## Explicit bindings

Sometimes it can be useful to refer to a parent node in the assertion tree.

This is possible to add an explicit binding at every node in the tree by appending the `->binding_name` to the key.

Given the input document:

```yaml
foo:
  bar: 4
  bat: 6
```

The following policy will compute a sum and bind the result to the `sum` binding. A descendant can then use `$sum` and use it:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
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
              (bar + bat)->sum:
                # get the `$sum` binding and compare it against `10`
                ($sum): 10
```

All binding are available to descendants, if a descendant creates a binding with a name that already exists the binding will be overridden for descendants only and it doesn't affect the bindings at upper levels in the tree.

In other words, a node in the tree always sees bindings that are defined in the parents and if a name is reused, the first binding with the given name wins when winding up the tree.

As a consequence, the policy below will evaluate to true:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar
      validate:
        assert:
          all:
          - foo:
              (bar + bat)->sum:
                ($sum + $sum)->sum:
                  ($sum): 20
                ($sum): 10
```

Finally, we can always access the current payload, policy and rule being evaluated using the built-in `$payload`, `$policy` and `$rule` bindings. No protection is made to prevent you from overriding those bindings though.


## Escaping projection

It can be necessary to prevent a projection under certain circumstances.

Consider the following document:

```yaml
foo:
  (bar): 4
  (baz):
  - 1
  - 2
  - 3
```

Here the `(bar)` key conflict with the projection syntax.
To workaround this situation, you can escape a projection by surrounding it with `\` characters like this:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar
      validate:
        assert:
          all:
          - foo:
              \(bar)\: 10
```

In this case, the leading and trailing `\` characters will be erased and the projection won't be applied.

Note that it's still possible to use the `~` modifier or to create a named binding with and escaped projection.

Keys like this are perfectly valid:

- `~index.\baz\`
- `\baz\@foo`
- `~index.\baz\@foo`
