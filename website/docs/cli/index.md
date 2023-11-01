# Overview

The `kyverno-json` Command Line Interface (CLI) can be used to:

* scan JSON or YAML files 
* launch a web application with a REST API
* launch a playground 

Here is an example of scanning an Terraform plan that creates an S3 bucket:

```sh
./kyverno-json scan --policy test/commands/scan/tf-s3/policy.yaml --payload test/commands/scan/tf-s3/payload.json
```

The output looks like:

```sh
Loading policies ...
Loading payload ...
Pre processing ...
Running ( evaluating 1 resource against 1 policy ) ...
- s3 / check-tags / (unknown) FAILED: all[0].check.planned_values.root_module.~.resources[0].values.(keys(tags_all)).(contains(@, 'Team')): Invalid value: false: Expected value: true
Done
```

## Installation

See [Install](../install.md) for the available options to install the CLI.

## Pre-processing payloads

You can provide preprocessing queries in [jmespath](https://jmespath.site) format to pre-process the input payload before evaluating *resources* against policies.

This is necessary if the input payload is not what you want to directly analyze.

For example, here is a partial JSON which was produced by converting a Terraform plan that creates an EC2 instance:

[kyverno/kyverno-json/main/test/commands/scan/tf-ec2/payload.json](https://github.com/kyverno/kyverno-json/blob/main/test/commands/scan/tf-ec2/payload.json)

```json
{
  "format_version": "1.2",
  "terraform_version": "1.5.7",
  "planned_values": {
    "root_module": {
      "resources": [
        {
          "address": "aws_instance.app_server",
          "mode": "managed",
          "type": "aws_instance",
          "name": "app_server",
          "provider_name": "registry.terraform.io/hashicorp/aws",
          "schema_version": 1,
          "values": {
            "ami": "ami-830c94e3",
            "credit_specification": [],
            "get_password_data": false,
            "hibernation": null,
            "instance_type": "t2.micro",
            "launch_template": [],
            "source_dest_check": true,
            "tags": {
              "Name": "ExampleAppServerInstance"
            },
            "tags_all": {
              "Name": "ExampleAppServerInstance"
            },
            "timeouts": null,
            "user_data_replace_on_change": false,
            "volume_tags": null
          },
   
          ...

```

To directly scan the `resources` element use `--pre-process planned_values.root_module.resources` as follows:

```sh
./kyverno-json scan --policy test/commands/scan/tf-ec2/policy.yaml --payload test/commands/scan/tf-ec2/payload.json --pre-process planned_values.root_module.resources
```

This command will produce the output:

```sh
Loading policies ...
Loading payload ...
Pre processing ...
Running ( evaluating 1 resource against 1 policy ) ...
- required-ec2-tags / require-team-tag / (unknown) PASSED
Done
```
