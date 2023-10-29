## kyverno-json docs

Generates reference documentation.

### Synopsis

Generates reference documentation.

The docs command generates CLI reference documentation.
It can be used to generate simple markdown files or markdown to be used for the website.

```
kyverno-json docs [flags]
```

### Examples

```
  # Generate simple markdown documentation
  kyverno-json docs -o . --autogenTag=false

  # Generate website documentation
  kyverno-json docs -o . --website

```

### Options

```
      --autogenTag      Determines if the generated docs should contain a timestamp (default true)
  -h, --help            help for docs
  -o, --output string   Output path (default ".")
      --website         Website version
```

### SEE ALSO

* [kyverno-json](kyverno-json.md)	 - kyverno-json is a CLI tool to apply policies to json resources.

