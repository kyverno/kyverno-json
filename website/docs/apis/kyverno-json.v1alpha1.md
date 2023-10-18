---
title: KyvernoJson (v1alpha1)
content_type: tool-reference
package: json.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the policy v1alpha1 API group</p>


## Resource Types 


- [ValidationPolicy](#json-kyverno-io-v1alpha1-ValidationPolicy)
- [ValidationPolicyList](#json-kyverno-io-v1alpha1-ValidationPolicyList)
  
## `ValidationPolicy`     {#json-kyverno-io-v1alpha1-ValidationPolicy}

**Appears in:**
    
- [ValidationPolicyList](#json-kyverno-io-v1alpha1-ValidationPolicyList)

<p>ValidationPolicy is the resource that contains the policy definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `json.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ValidationPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`PolicySpec`](#json-kyverno-io-v1alpha1-PolicySpec) | :white_check_mark: |  | <p>Policy spec.</p> |

## `ValidationPolicyList`     {#json-kyverno-io-v1alpha1-ValidationPolicyList}

<p>ValidationPolicyList is a list of Policy instances.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `json.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ValidationPolicyList` |
| `metadata` | [`meta/v1.ListMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `items` | [`[]ValidationPolicy`](#json-kyverno-io-v1alpha1-ValidationPolicy) | :white_check_mark: |  | *No description provided.* |

## `Any`     {#json-kyverno-io-v1alpha1-Any}

**Appears in:**
    
- [Assertion](#json-kyverno-io-v1alpha1-Assertion)
- [ContextEntry](#json-kyverno-io-v1alpha1-ContextEntry)
- [Match](#json-kyverno-io-v1alpha1-Match)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `Value` | `interface{}` | :white_check_mark: | :white_check_mark: | *No description provided.* |

## `Assert`     {#json-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [ValidationRule](#json-kyverno-io-v1alpha1-ValidationRule)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `any` | [`[]Assertion`](#json-kyverno-io-v1alpha1-Assertion) | :white_check_mark: |  | <p>Any allows specifying resources which will be ORed.</p> |
| `all` | [`[]Assertion`](#json-kyverno-io-v1alpha1-Assertion) | :white_check_mark: |  | <p>All allows specifying resources which will be ANDed.</p> |

## `Assertion`     {#json-kyverno-io-v1alpha1-Assertion}

**Appears in:**
    
- [Assert](#json-kyverno-io-v1alpha1-Assert)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `message` | `string` | :white_check_mark: |  | <p>Message is the variable associated message.</p> |
| `check` | [`Any`](#json-kyverno-io-v1alpha1-Any) | :white_check_mark: |  | <p>Check is the assertion check definition.</p> |

## `ContextEntry`     {#json-kyverno-io-v1alpha1-ContextEntry}

**Appears in:**
    
- [ValidationRule](#json-kyverno-io-v1alpha1-ValidationRule)

<p>ContextEntry adds variables and data sources to a rule Context.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the variable name.</p> |
| `variable` | [`Any`](#json-kyverno-io-v1alpha1-Any) | :white_check_mark: |  | <p>Variable defines an arbitrary JMESPath context variable that can be defined inline.</p> |

## `Match`     {#json-kyverno-io-v1alpha1-Match}

**Appears in:**
    
- [ValidationRule](#json-kyverno-io-v1alpha1-ValidationRule)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `any` | [`[]Any`](#json-kyverno-io-v1alpha1-Any) | :white_check_mark: |  | <p>Any allows specifying resources which will be ORed.</p> |
| `all` | [`[]Any`](#json-kyverno-io-v1alpha1-Any) | :white_check_mark: |  | <p>All allows specifying resources which will be ANDed.</p> |

## `PolicySpec`     {#json-kyverno-io-v1alpha1-PolicySpec}

**Appears in:**
    
- [ValidationPolicy](#json-kyverno-io-v1alpha1-ValidationPolicy)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `rules` | [`[]ValidationRule`](#json-kyverno-io-v1alpha1-ValidationRule) | :white_check_mark: |  | <p>Rules is a list of Rule instances. A Policy contains multiple rules and each rule can validate, mutate, or generate resources.</p> |

## `ValidationRule`     {#json-kyverno-io-v1alpha1-ValidationRule}

**Appears in:**
    
- [PolicySpec](#json-kyverno-io-v1alpha1-PolicySpec)

| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is a label to identify the rule, It must be unique within the policy.</p> |
| `context` | [`[]ContextEntry`](#json-kyverno-io-v1alpha1-ContextEntry) | :white_check_mark: |  | <p>Context defines variables and data sources that can be used during rule execution.</p> |
| `match` | [`Match`](#json-kyverno-io-v1alpha1-Match) | :white_check_mark: |  | <p>Match defines when this policy rule should be applied. The match criteria can include resource information (e.g. kind, name, namespace, labels) and admission review request information like the user name or role. At least one kind is required.</p> |
| `exclude` | [`Match`](#json-kyverno-io-v1alpha1-Match) | :white_check_mark: |  | <p>Exclude defines when this policy rule should not be applied. The exclude criteria can include resource information (e.g. kind, name, namespace, labels) and admission review request information like the name or role.</p> |
| `assert` | [`Assert`](#json-kyverno-io-v1alpha1-Assert) | :white_check_mark: |  | <p>Assert is used to validate matching resources.</p> |

  