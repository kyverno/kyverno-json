# Overview

`kyverno-json` uses [JMESPath community edition](https://jmespath.site/), a modern JMESPath implementation with lexical scopes support.

The current resource, policy and rule are always available using the following builtin bindings:

| Binding | Usage |
|---|---|
| `$resource` | Current resource being analysed |
| `$policy` | Current policy being executed |
| `$rule` | Current rule being evaluated |

No protection is made to prevent you from overriding those bindings though.
