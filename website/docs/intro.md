# Introduction

`kyverno-json` is a CLI tool very similar to the [Kyverno CLI](https://github.com/kyverno/kyverno/tree/main/cmd/cli/kubectl-kyverno).

The difference is that `kyverno-json` can apply policies to abitrary json or yaml payloads.

Policy definition syntax looks a lot like the [Kyverno policy](https://kyverno.io/docs/kyverno-policies/) definition syntax but is more generic and flexible.

This was needed to allow working with arbitrary payloads, not just [Kubernetes](https://kubernetes.io) ones.

## Pre-processing

Additionally, you can provide preprocessing queries in [jmespath](https://jmespath.site) format to preprocess the input payload before evaluating *resources* against policies.

This is necessary if the input payload is not what you want to directly analyse.
