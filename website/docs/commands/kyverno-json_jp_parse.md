## kyverno-json jp parse

Parses jmespath expression and prints corresponding AST.

### Synopsis

Parses jmespath expression and prints corresponding AST.


```
kyverno-json jp parse [-f file|expression]... [flags]
```

### Examples

```
  # Parse expression
  kyverno-json jp parse 'request.object.metadata.name | truncate(@, `9`)'

  # Parse expression from a file
  kyverno-json jp parse -f my-file

  # Parse expression from stdin
  kyverno-json jp parse

  # Parse multiple expressionxs
  kyverno-json jp parse -f my-file1 -f my-file-2 'request.object.metadata.name | truncate(@, `9`)'

```

### Options

```
  -f, --file strings   Read input from a JSON or YAML file instead of stdin
  -h, --help           help for parse
```

### SEE ALSO

* [kyverno-json jp](kyverno-json_jp.md)	 - Provides a command-line interface to JMESPath, enhanced with custom functions.

