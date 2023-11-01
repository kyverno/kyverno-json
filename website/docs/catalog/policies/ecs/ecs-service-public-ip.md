---
tags:
- aws
- aws/ecs
---
# ECS public IP

## Description

This Policy ensures that ECS services do not have public IP addresses assigned to them automatically.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-service-public-ip.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-service-public-ip.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/ecs/ecs-service-public-ip.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-service-public-ip.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that ECS services do not have
      public IP addresses assigned to them automatically.
    title.policy.kyverno.io: ECS public IP
  creationTimestamp: null
  labels:
    ecs.aws.tags.kyverno.io: ecs-service
  name: ecs-public-ip
spec:
  rules:
  - assert:
      all:
      - check:
          values:
            ~.network_configuration:
              (contains('$allowed-values', @.assign_public_ip)): false
        message: ECS services should not have public IP addresses assigned to them
          automatically
    context:
    - name: allowed-values
      variable:
      - false
    match:
      any:
      - type: aws_ecs_service
    name: ecs-public-ip
```
