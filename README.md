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
./json-kyverno --payload ./tf.plan.json --pre-process "planned_values.root_module.resources" --policy ./policy.yaml
```
