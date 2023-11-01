---
tags:
- aws
- aws/ecs
---
# ECS cluster enable logging

## Description

This Policy ensures that ECS clusters have logging enabled.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-cluster-enable-logging.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-cluster-enable-logging.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/ecs/ecs-cluster-enable-logging.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-cluster-enable-logging.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that ECS clusters have logging
      enabled.
    title.policy.kyverno.io: ECS cluster enable logging
  creationTimestamp: null
  labels:
    ecs.aws.tags.kyverno.io: ecs-cluster
  name: ecs-cluster-enable-logging
spec:
  rules:
  - assert:
      all:
      - check:
          values:
            ~.configuration:
              ~.execute_command_configuration:
                (contains($forbidden_values, @.logging)): false
        message: ECS Cluster should enable logging of ECS Exec
    context:
    - name: forbidden_values
      variable:
      - NONE
    match:
      any:
      - type: aws_ecs_cluster
    name: ecs-cluster-enable-logging
```
