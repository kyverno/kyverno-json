# Test for running the api

## Create a cluster

```bash
make kind-create
```

## Install CRDs

```bash
make install-crds
```

## Deploy a policy

```bash
kubectl apply -f - <<EOF
apiVersion: json.kyverno.io/v1alpha1
kind: ValidatingPolicy
metadata:
  name: test
spec:
  rules:
    - name: foo-bar-4
      assert:
        all:
        - check:
            foo:
              bar: 4
EOF
```

## Run KyvernoJson

```bash
make build

./kyverno-json serve
```

## Call the KyvernoJson API

```bash
curl -X POST http://localhost:8080/api/scan -H 'Content-Type: application/json' -d @- <<EOF
{
    "payload": {
        "foo": {
            "bar": 4
        }
    }
}
EOF
```
