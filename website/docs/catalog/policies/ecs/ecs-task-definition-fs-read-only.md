---
tags:
- aws
- aws/ecs
---
# ECS require filesystem read only

## Description

This Policy ensures that ECS Fargate services runs on the latest Fargate platform version.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-task-definition-fs-read-only.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-task-definition-fs-read-only.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/ecs/ecs-task-definition-fs-read-only.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-task-definition-fs-read-only.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that ECS Fargate services runs
      on the latest Fargate platform version.
    title.policy.kyverno.io: ECS require filesystem read only
  creationTimestamp: null
  labels:
    ecs.aws.tags.kyverno.io: ecs-task-definition
  name: fs-read-only
spec:
  rules:
  - assert:
      any:
      - check:
          values:
            ~.(json_parse(container_definitions)):
              readonlyRootFilesystem: true
        message: ECS containers should only have read-only access to root filesystems
    match:
      any:
      - type: aws_ecs_task_definition
    name: require-fs-read-only
```
