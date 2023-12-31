## kyverno-json jp

Provides a command-line interface to JMESPath, enhanced with custom functions.

### Synopsis

Provides a command-line interface to JMESPath, enhanced with custom functions.


```
kyverno-json jp [flags]
```

### Examples

```
  # List functions
  kyverno-json jp function

  # Evaluate query
  kyverno-json jp query -i object.yaml 'request.object.metadata.name | truncate(@, `9`)'

  # Parse expression
  kyverno-json jp parse 'request.object.metadata.name | truncate(@, `9`)'

```

### Options

```
  -h, --help   help for jp
```

### SEE ALSO

* [kyverno-json](kyverno-json.md)	 - kyverno-json is a CLI tool to apply policies to json resources.
* [kyverno-json jp function](kyverno-json_jp_function.md)	 - Provides function informations.
* [kyverno-json jp parse](kyverno-json_jp_parse.md)	 - Parses jmespath expression and prints corresponding AST.
* [kyverno-json jp query](kyverno-json_jp_query.md)	 - Provides a command-line interface to JMESPath, enhanced with Kyverno specific custom functions.

