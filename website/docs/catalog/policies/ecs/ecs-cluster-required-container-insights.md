---
tags:
- aws
- aws/ecs
---
# ECS requires container insights

## Description

This Policy ensures that ECS clusters have container insights enabled.

## Install

### In cluster

```bash
kubectl apply -f https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-cluster-required-container-insights.yaml
```

### Download locally

```bash
curl -O https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-cluster-required-container-insights.yaml
```

## Manifest

[Original policy](https://github.com/kyverno/kyverno-json/blob/main/catalog/ecs/ecs-cluster-required-container-insights.yaml)
[Raw](https://raw.githubusercontent.com/kyverno/kyverno-json/main/catalog/ecs/ecs-cluster-required-container-insights.yaml)

```yaml
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  annotations:
    description.policy.kyverno.io: This Policy ensures that ECS clusters have container
      insights enabled.
    title.policy.kyverno.io: ECS requires container insights
  creationTimestamp: null
  labels:
    ecs.aws.tags.kyverno.io: ecs-cluster
  name: required-container-insights
spec:
  rules:
  - assert:
      all:
      - check:
          values:
            ~.setting:
              name: containerInsights
              value: enabled
        message: Container insights should be enabled on ECS cluster
    match:
      any:
      - type: aws_ecs_cluster
    name: required-container-insights
```
