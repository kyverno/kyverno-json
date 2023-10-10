# Apply policies on a Dockerfile

1. Download a Dockerfile

```
curl https://raw.githubusercontent.com/nirmata/kyverno-notation-aws/main/Dockerfile /tmp/Dockefile
```

2. Convert to JSON

Install `dockerfile-json`: https://github.com/keilerkonzept/dockerfile-json#get-it

```
dockerfile-json ~/go/src/github.com/jimbugwadia/kyverno-notation-aws/Dockerfile | jq > input.json
```

3. Apply policy

```
/kyverno-json scan --policy testdata/dockerfile/policy-check-external.yaml --payload testdata/dockerfile/input.json
```

Results:

```
Loading policies ...
Loading payload ...
Pre processing ...
Running ( evaluating 1 resource against 1 policy ) ...
- check-dockerfile / no-http /  FAILED: HTTP calls are not allowed
- check-dockerfile / no-curl-wget /  FAILED: curl / wget are not allowed
Done
```