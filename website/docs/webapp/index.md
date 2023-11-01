# Usage

`kyverno-json` can be deployed as a web application with a REST API. This is useful for deployments when a long running service that processes policy requests is desired.

## Managing Policies

With `kyverno-json` policies are managed as Kubernetes resources. This means that you can use Kubernetes APIs, `kubectl`, GitOps, or any other Kubernetes management tool to manage policies.

## Usage

Here is a complete demonstration of how to use `kyverno-json` as an web application:

**Install CRDs**

Install the CRD for `kyverno-json`:

```sh
kubectl apply -f .crds/json.kyverno.io_validatingpolicies.yaml
```

**Install policies:**

Install a sample policy:

```sh
kubectl apply -f test/commands/scan/dockerfile/policy.yaml
```

**Prepare the payload**

The payload is a JSON object with two fields:

| Name            | Type             | Required     |
| --------------- | ---------------- | ------------ |
| `payload`       | Object           | Y            |
| `preprocessors` | Array of Strings | N            |


You can construct a sample payload for the Dockerfile policy using:

```sh
cat test/commands/scan/dockerfile/payload.json | jq '{"payload": .}' > /tmp/webapp-payload.json
```

Run the web application

```sh
./kyverno-json serve
```

This will show the output:

```sh
2023/10/29 23:46:11 configured route /api/scan
2023/10/29 23:46:11 listening to requests on 0.0.0.0:8080
```

Send the REST API request

```sh
curl http://localhost:8080/api/scan -X POST -H "Content-Type: application/json" -d @/tmp/webapp-payload.json | jq
```

The configured policies will be applied to the payload and the results will be returned back:

```sh
{
  "results": [
    {
      "policy": "check-dockerfile",
      "rule": "deny-external-calls",
      "status": "fail",
      "message": "HTTP calls are not allowed: all[0].check.~.(Stages[].Commands[].Args[].Value)[0].(contains(@, 'https://') || contains(@, 'http://')): Invalid value: true: Expected value: false; wget is not allowed: all[3].check.~.(Stages[].Commands[].CmdLine[])[0].(contains(@, 'wget')): Invalid value: true: Expected value: false"
    }
  ]
}
```

## Helm Chart

The web application can be installed and managed in a Kubernetes cluster using Helm. 

See details at: https://github.com/kyverno/kyverno-json/tree/main/charts/kyverno-json
