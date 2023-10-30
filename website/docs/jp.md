# Overview

`kyverno-json` uses [JMESPath community edition](https://jmespath.site/), a modern JMESPath implementation with lexical scopes support.

The current *payload*, *policy* and *rule* are always available using the following builtin bindings:

| Binding | Usage |
|---|---|
| `$payload` | Current payload being analysed |
| `$policy` | Current policy being executed |
| `$rule` | Current rule being evaluated |

!!! warning

    No protection is made to prevent you from overriding those bindings.
