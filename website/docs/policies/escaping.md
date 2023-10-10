# Escaping projection

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
              \(bar)\: 10
```

In this case, the leading and trailing `\` characters will be erased and the projection won't be applied.

Note that it's still possible to use the `~` modifier or to create a named binding with and escaped projection.

Keys like this are perfectly valid:

- `~index.\baz\`
- `\baz\@foo`
- `~index.\baz\@foo`
