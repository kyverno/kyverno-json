# json-kyverno

## create a terraform plan in json format

```console
# init terraform
terraform init

# create a plan
terraform plan -out=tf.plan

# show plan in json
terraform show -json tf.plan > tf.plan.json
```

## build json-kyverno

```console
make build
```

## invoke json-kyverno

```console
./json-kyverno --plan ./tf.plan.json --policy ./policy.yaml
```
