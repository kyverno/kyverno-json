# Introduction

`kyverno-json` extends Kyverno policies to perform simple and efficient validation of data in JSON or YAML format. With `kyverno-json`, you can now use Kyverno policies to validate:

- Terraform files
- Dockerfiles
- Cloud configurations
- Authorization requests

Simply convert your runtime or configuration data to JSON, and use Kyverno to audit or enforce policies for security and best practices compliance.

`kyverno-json` can be run as a:

1. [A Command Line Interface (CLI)](./cli/index.md)
2. [A web application with a REST API](./webapp/index.md)
3. [A Golang library](./go-library/index.md)
