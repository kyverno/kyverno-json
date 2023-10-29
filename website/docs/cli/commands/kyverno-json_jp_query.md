## kyverno-json jp query

Provides a command-line interface to JMESPath, enhanced with Kyverno specific custom functions.

### Synopsis

Provides a command-line interface to JMESPath, enhanced with Kyverno specific custom functions.


```
kyverno-json jp query [-i input] [-q query|query]... [flags]
```

### Examples

```
  # Evaluate query
  kyverno-json jp query -i object.yaml 'request.object.metadata.name | truncate(@, `9`)'

  # Evaluate query
  kyverno-json jp query -i object.yaml -q query-file

  # Evaluate multiple queries
  kyverno-json jp query -i object.yaml -q query-file-1 -q query-file-2 'request.object.metadata.name | truncate(@, `9`)'

```

### Options

```
  -c, --compact         Produce compact JSON output that omits non essential whitespace
  -h, --help            help for query
  -i, --input string    Read input from a JSON or YAML file instead of stdin
  -q, --query strings   Read JMESPath expression from the specified file
  -u, --unquoted        If the final result is a string, it will be printed without quotes
```

### SEE ALSO

* [kyverno-json jp](kyverno-json_jp.md)	 - Provides a command-line interface to JMESPath, enhanced with custom functions.

