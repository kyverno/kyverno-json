# tf-kyverno

## create a terraform plan in json format

```console
# create a plan
terraform plan -out=tf.plan

# show plan in json
terraform show -json tf.plan > tf.plan.json
```

## invoke tf-kyverno

```console
./tf-kyverno \
    --plan ./tf.plan.json \
    --policy ./policy.yaml
```