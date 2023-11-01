---
tags:
- aws
- aws/ecs
---
# ECS require latest platform fargate

## Description

This Policy ensures that ECS Fargate services runs on the latest Fargate platform version.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-service-required-latest-platform-fargate.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-service-required-latest-platform-fargate.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/ecs/ecs-service-required-latest-platform-fargate.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-service-required-latest-platform-fargate.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that ECS Fargate services runs
      on the latest Fargate platform version.
    title.policy.kyverno.io: ECS require latest platform fargate
  creationTimestamp: null
  labels:
    ecs.aws.tags.kyverno.io: ecs-service
  name: required-latest-platform-fargate
spec:
  rules:
  - assert:
      all:
      - check:
          values:
            platform_version: LATEST
        message: ECS Fargate services should run on the latest Fargate platform version
    context:
    - name: pv
      variable: platform_version
    match:
      any:
      - type: aws_ecs_service
        values:
          launch_type: FARGATE
    name: required-latest-platform
```
