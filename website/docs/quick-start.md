# Quick Start

## Validate a Terraform Plan

In this example we will use a Kyverno policy to validate a Terraform plan:

### Create the payload

Here is a Terraform plan that creates an AWS S3 bucket:

```terraform
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-west-2"
}

resource "aws_s3_bucket" "example" {
  bucket = "my-tf-test-bucket"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}
```

You can convert this to JSON using the following commands:

*output the plan:*
```sh
terraform plan -out tfplan.binary
```
*convert to JSON:*
```sh
terraform show -json tfplan.binary | jq > payload.json
```

### Create the policy

Create a `policy.yaml` file and paste the content below that checks for required labels:


```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: s3
spec:
  rules:
    - name: check-tags
      assert:
        all:
        - check:
            planned_values:
              root_module:
                ~.resources:
                  values:
                    (keys(tags_all)):
                      (contains(@, 'Environment')): true
                      (contains(@, 'Name')): true
                      (contains(@, 'Team')): true
```

### Scan the payload

With the payload and policy above, we can invoke `kyverno-json` with the command below:

```bash
kyverno-json scan --payload payload.json --policy policy.yaml
```

The plan shown above will fail as it does not contain the `Team` tag.

```sh
Loading policies ...
Loading payload ...
Pre processing ...
Running ( evaluating 1 resource against 1 policy ) ...
- s3 / check-tags / (unknown) FAILED: all[0].check.planned_values.root_module.~.resources[0].values.(keys(tags_all)).(contains(@, 'Team')): Invalid value: false: Expected value: true
Done
```

## Validate a Kubernetes Resource

For this example we will use a [Kubernetes](https://kubernetes.io) `Pod` payload.

### Create the payload

Create a `payload.yaml` file and paste the Pod declaration below in it:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pods-simple-pod
spec:
  containers:
    - command:
        - sleep
        - "3600"
      image: busybox:latest
      name: pods-simple-container
```

This is a simple `Pod` with one container running the `busybox` latest docker image.

Using the `latest` tag of an image is a bad practice. Let's write a policy to detect this.

### Create the policy

Create a `policy.yaml` file and paste the content below to block `latest` images:

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: pod-policy
spec:
  rules:
    - name: no-latest
      # Match payloads corresponding to pods
      match:
        any:
        - apiVersion: v1
          kind: Pod
      validate:
        message: Pod `{{ metadata.name }}` uses an image with tag `latest`
        assert:
          all:
          - spec:
              # Iterate over pod containers
              # Note the `~.` modifier, it means we want to iterate over array elements in descendants
              ~.containers:
                image:
                  # Check that an image tag is present
                  (contains(@, ':')): true
                  # Check that the image tag is not `:latest`
                  (ends_with(@, ':latest')): false
```

This policy iterates over pod containers, checking that the container image has a tag specified and that the tag being used is not `latest`.

### Scan the payload

With the payload and policy above, we can invoke `kyverno-json` with the command below:

```bash
kyverno-json scan --payload payload.yaml --policy policy.yaml
```

This produces the output:

```bash
Loading policies ...
Loading payload ...
Pre processing ...
Running ( evaluating 1 resource against 1 policy ) ...
- pod-policy / no-latest /  FAILED: Pod `pods-simple-pod` uses an image with tag `latest`
Done
```
