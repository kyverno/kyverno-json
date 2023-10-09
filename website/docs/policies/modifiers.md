# Projection modifiers

Assertion tree expressions support modifiers to influence the way projected values are processed.

The `~` modifier applies to arrays and maps, it mean the input array or map elements will be processed individually by descendants.
When the `~` modifier is not used, descendants receive the whole array, not individual elements.

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
              # the content of the `bar` field will be compared against `[1, 2, 3]`
              bar:
              - 1
              - 2
              - 3
```

With the `~` modifier, we can apply descendants to all elements in the array individually.
The policy below ensures that all elements in the input array are `< 5`:

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
              # with the `~` modifier all elements in the `[1, 2, 3]` array are processed individually and passed to descendants
              ~.bar:
                # the expression `(@ < `5`)` is evaluated for every element and the result is expected to be `true`
                (@ < `5`): true
```

The `~` modifier supports binding the index of the element being processed to a named binding with the following syntax `~index_name.bar`. When this is used, we can access the element index in descendants with `$index_name`.

When used with a map, the named binding receives the key of the element being processed.
