# Playground examples docs

This docs contains information to manage playground examples.

## Modify playground examples

To add, update or remove a playground example edit the [playground-examples.yaml](../playground-examples.yaml) file.

This file contains two nested maps, first level is the example category and second level is example name:

```yaml
example category 1:
    example name 1:
        policy: path/to/policy/file (yaml or json)
        payload: path/to/payload/file (yaml or json)
    example name 2:
        policy: path/to/policy/file (yaml or json)
        payload: path/to/payload/file (yaml or json)
example category 2:
    # ...
```

Once the file edited, run `make codegen-playground-examples` to update the [data.json](../website/playground/assets/data.json) file.
