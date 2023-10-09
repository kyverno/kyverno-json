# kyverno-json

This CLI tool is very similar to the [Kyverno CLI](https://github.com/kyverno/kyverno/tree/main/cmd/cli/kubectl-kyverno) tool.

The difference is that this CLI tool can apply policies to abitrary json or yaml payloads.

Policy definition syntax is looks a lot like the [Kyverno policy](https://kyverno.io/docs/kyverno-policies/) definition syntax but is more generic and flexible.
This was needed to allow working with arbitrary payloads, not just [Kubernetes](https://kubernetes.io) ones.
Those differences are detailed in the [section below](#differences-with-with-kyverno-policy-definition-syntax).

Additionally, you can provide preprocessing queries in [jmespath](https://jmespath.site) format to preprocess the input payload before evaluating *resources* against policies.
This is necessary if the input payload is not what you want to directly analyse.
Preprocessing is detailed in the following [section](#preprocessing).

## Differences with with Kyverno policy definition syntax

Sections below highlight the main differences between polcies used by this tool and [Kyverno policies](https://kyverno.io/docs/kyverno-policies/).

### Different `apiVersion` and `kind`

Both [Kyverno policies](https://kyverno.io/docs/kyverno-policies/) and policies used by this tool are defined using [Kubernetes](https://kubernetes.io) manifests.

They don't use the same `apiVersion` and `kind` though.

[Kyverno policies](https://kyverno.io/docs/kyverno-policies/) belong to the `kyverno.io` group, exist in multiple versions (`v1`, `v2beta1`) and can be of kind `Policy` or `ClusterPolicy`.

Policies for this tool belong to the `json.kyverno.io` group, exist only in `v1alpha1` version and can only be of kind `Policy`.

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

The concept of clustered vs namespaced resources exist only in the [Kubernetes](https://kubernetes.io) world and it didn't make sense to reproduce the same pattern in this tool.

### Different `match` and `exclude` statements

Both [Kyverno policies](https://kyverno.io/docs/kyverno-policies/) and policies used by this tool can match and exclude *resources* when being evaluated.

[Kyverno policies](https://kyverno.io/docs/kyverno-policies/) use [Kubernetes](https://kubernetes.io) specific constructs for that matter that didn't map well with arbitrary payloads.

This tool uses [assertion trees](#assertion-trees-replace-pattern-matching) to implement `match` and `exclude` statements:

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

### Different `jmesPath` implementation

This tool uses [jmespath-community/go-jmespath](https://github.com/jmespath-community/go-jmespath), a more modern implementation than the one used in [Kyverno](https://kyverno.io).

This implementation supports the `let` feature and this tool leverages it to implement context entries:

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
      - name: expectedTeam
        variable: Kyverno
      validate:
        message: Bucket `{{ name }}` ({{ address }}) does not have the required Team tag {{ $expectedTeam }}
        assert:
          all:
          - values:
              tags:
                Team: ($expectedTeam)
```

Note that all context entries are lazily evaluated, a context entry will only be evaluated once. They can be used in all [assertion trees](#assertion-trees-replace-pattern-matching), including `match` and `exclude` statements.

### No preconditions, pattern operators, anchors or wildcards

Policies used by this tool don't support `preconditions`, pattern operators, anchors or wildcards.

Most of the time `preconditions` can be replaced by the more flexible `match` and `exclude` statements.

Pattern operators, anchors and wildcards can be replaced with an improved pattern matching system.
The new pattern matching system is called *assertion trees*, this is detailed [below](#assertion-trees-replace-pattern-matching).

### Assertion trees replace pattern matching

[Kyverno policies](https://kyverno.io/docs/kyverno-policies/) started with a declarative approach but slowly adopted the imperative approach too, because of the limitations in the implemented declarative approach.

This tool tries to be as declarative as possible, for now `forEach`, pattern operators, anchors and wildcards are not supported are not supported.
Hopefully we won't need to adopt an imperative approach anymore.

Instead, assertion trees can now be used to express complex and dynamic conditions by using [jmespath](https://jmespath.site) expressions. Those expressions represent projections of the being analysed *resource* and the result of this projection is passed to descendants for further analysis.

All comparisons happen in the leaves of the assertion tree.

Given the input payload below:

```yaml
foo:
  baz: true
  bar: 4
  bat: 6
```

It is now possible to write a validation tree like this:

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

#### Projection modifiers

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

#### Explicit bindings

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

Finally, we can always access the current resource, policy and rule being evaluated using the builtin `$resource`, `$policy` and `$rule` bindings. No protection is made to prevent you from overriding those bindings though.

#### Escaping projection

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

Here the `(bar)` key conflict with the projection syntax used.
To workaround this issue, you can escape a projection by surrounding it with `/` characters like this:

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
              /(bar)/: 10
```

In this case, the leading and trailing `/` characters will be erased and the projection won't be applied.

Note that it's still possible to use the `~` modifier or to create a named binding with and escaped projection.

Keys like this are perfectly valid:
- `~index./baz/`
- `/baz/@foo`
- `~index./baz/@foo`

## SDK

This CLI tool contains an initial implementation of an SDK to allow flexible creation of dedicated policy engines.

The [json-engine](./pkg/json-engine/) at the heart of this tool is built by assembling blocks provided by the [engine](./pkg/engine/) SDK.

## Build kyverno-json

To build this tool locally, simply run:

```console
make build
```

## Preprocessing

When the input payload is not what you want to analyse directly, you can provide one or more [jmespath](https://jmespath.site) expressions to preprocess the data.
The policies will be evaluated against the result of the preprocessing step.

Traditionnally, policies apply to *resources* and this is how this tool implements policy evaluation:

```
loop through all resources {
    loop through all policies {
        loop through all rules in the policy {
            evaluate the resource against the rule
        }
    }
}
```

Note that if you provide a single payload, the tool will internally wrap it in an array of one element.

So imagine an input payload similar to this:

```yaml
version: 1.2.3
creationDate: '2023-09-29'
resources:
- type: something
  name: foo
  spec:
    # ...
- type: something else
  name: bar
  spec:
    # ...
```

The *resources* you want to analyse are located under the `resources` stanza, and your policies are probably written to work on those *resources*.
In order to extract the data under the `resources` stanza before processing happens you can specify the `--pre-process "resources"` when invoking the tool.

You can chain mutliple preprocessing queries by specifying the `--pre-process` flag multiple times.
There is no limitation in a preprocessing [jmespath](https://jmespath.site) expression.

## Invoke kyverno-json

```console
# with yaml payload
./kyverno-json scan --payload ./testdata/foo-bar/payload.yaml --policy ./testdata/foo-bar/policy.yaml

# with json payload (and pre processing)
./kyverno-json scan --payload ./testdata/tf-plan/tf.plan.json --pre-process "planned_values.root_module.resources" --policy ./testdata/tf-plan/policy.yaml
```

## Documentation

- User documentation can be found in [docs/user](./docs/user/README.md)
- Dev documentation can be found in [docs/dev](./docs/dev/README.md)
