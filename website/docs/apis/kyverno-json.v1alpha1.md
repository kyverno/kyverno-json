---
title: KyvernoJson (v1alpha1)
content_type: tool-reference
package: json.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the policy v1alpha1 API group</p>


## Resource Types 


- [ValidatingPolicy](#json-kyverno-io-v1alpha1-ValidatingPolicy)
- [ValidatingPolicyList](#json-kyverno-io-v1alpha1-ValidatingPolicyList)
  
## `ValidatingPolicy`     {#json-kyverno-io-v1alpha1-ValidatingPolicy}

**Appears in:**
    
- [ValidatingPolicyList](#json-kyverno-io-v1alpha1-ValidatingPolicyList)

<p>ValidatingPolicy is the resource that contains the policy definition.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `json.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ValidatingPolicy` |
| `metadata` | [`meta/v1.ObjectMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#objectmeta-v1-meta) |  |  | <p>Standard object's metadata.</p> |
| `spec` | [`ValidatingPolicySpec`](#json-kyverno-io-v1alpha1-ValidatingPolicySpec) | :white_check_mark: |  | <p>Policy spec.</p> |

## `ValidatingPolicyList`     {#json-kyverno-io-v1alpha1-ValidatingPolicyList}

<p>ValidatingPolicyList is a list of ValidatingPolicy instances.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `apiVersion` | `string` | :white_check_mark: | | `json.kyverno.io/v1alpha1` |
| `kind` | `string` | :white_check_mark: | | `ValidatingPolicyList` |
| `metadata` | [`meta/v1.ListMeta`](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#listmeta-v1-meta) | :white_check_mark: |  | *No description provided.* |
| `items` | [`[]ValidatingPolicy`](#json-kyverno-io-v1alpha1-ValidatingPolicy) | :white_check_mark: |  | *No description provided.* |

## `Any`     {#json-kyverno-io-v1alpha1-Any}

**Appears in:**
    
- [Assertion](#json-kyverno-io-v1alpha1-Assertion)
- [ContextEntry](#json-kyverno-io-v1alpha1-ContextEntry)
- [Match](#json-kyverno-io-v1alpha1-Match)

<p>Any can be any type.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `Value` | `any` |  | :white_check_mark: | <p>Value contains the value of the Any object.</p> |

## `Assert`     {#json-kyverno-io-v1alpha1-Assert}

**Appears in:**
    
- [ValidatingRule](#json-kyverno-io-v1alpha1-ValidatingRule)

<p>Assert defines collections of assertions.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `any` | [`[]Assertion`](#json-kyverno-io-v1alpha1-Assertion) |  |  | <p>Any allows specifying assertions which will be ORed.</p> |
| `all` | [`[]Assertion`](#json-kyverno-io-v1alpha1-Assertion) |  |  | <p>All allows specifying assertions which will be ANDed.</p> |

## `Assertion`     {#json-kyverno-io-v1alpha1-Assertion}

**Appears in:**
    
- [Assert](#json-kyverno-io-v1alpha1-Assert)

<p>Assertion contains an assertion tree associated with a message.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `message` | `string` |  |  | <p>Message is the message associated message.</p> |
| `check` | [`Any`](#json-kyverno-io-v1alpha1-Any) | :white_check_mark: |  | <p>Check is the assertion check definition.</p> |

## `ContextEntry`     {#json-kyverno-io-v1alpha1-ContextEntry}

**Appears in:**
    
- [ValidatingRule](#json-kyverno-io-v1alpha1-ValidatingRule)

<p>ContextEntry adds variables and data sources to a rule context.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the entry name.</p> |
| `variable` | [`Any`](#json-kyverno-io-v1alpha1-Any) |  |  | <p>Variable defines an arbitrary variable.</p> |

## `Feedback`     {#json-kyverno-io-v1alpha1-Feedback}

**Appears in:**
    
- [ValidatingRule](#json-kyverno-io-v1alpha1-ValidatingRule)

<p>Feedback contains a feedback entry.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is the feedback entry name.</p> |
| `value` | `string` | :white_check_mark: |  | <p>Value is the feedback entry value (a JMESPath expression).</p> |

## `Match`     {#json-kyverno-io-v1alpha1-Match}

**Appears in:**
    
- [ValidatingRule](#json-kyverno-io-v1alpha1-ValidatingRule)

<p>Match defines collections of assertion trees.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `any` | [`[]Any`](#json-kyverno-io-v1alpha1-Any) |  |  | <p>Any allows specifying assertion trees which will be ORed.</p> |
| `all` | [`[]Any`](#json-kyverno-io-v1alpha1-Any) |  |  | <p>All allows specifying assertion trees which will be ANDed.</p> |

## `ValidatingPolicySpec`     {#json-kyverno-io-v1alpha1-ValidatingPolicySpec}

**Appears in:**
    
- [ValidatingPolicy](#json-kyverno-io-v1alpha1-ValidatingPolicy)

<p>ValidatingPolicySpec contains the policy spec.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `rules` | [`[]ValidatingRule`](#json-kyverno-io-v1alpha1-ValidatingRule) | :white_check_mark: |  | <p>Rules is a list of ValidatingRule instances.</p> |

## `ValidatingRule`     {#json-kyverno-io-v1alpha1-ValidatingRule}

**Appears in:**
    
- [ValidatingPolicySpec](#json-kyverno-io-v1alpha1-ValidatingPolicySpec)

<p>ValidatingRule defines a validating rule.</p>


| Field | Type | Required | Inline | Description |
|---|---|---|---|---|
| `name` | `string` | :white_check_mark: |  | <p>Name is a label to identify the rule, It must be unique within the policy.</p> |
| `context` | [`[]ContextEntry`](#json-kyverno-io-v1alpha1-ContextEntry) |  |  | <p>Context defines variables and data sources that can be used during rule execution.</p> |
| `match` | [`Match`](#json-kyverno-io-v1alpha1-Match) |  |  | <p>Match defines when this policy rule should be applied.</p> |
| `exclude` | [`Match`](#json-kyverno-io-v1alpha1-Match) |  |  | <p>Exclude defines when this policy rule should not be applied.</p> |
| `identifier` | `string` |  |  | <p>Identifier declares a JMESPath expression to extract a name from the payload.</p> |
| `feedback` | [`[]Feedback`](#json-kyverno-io-v1alpha1-Feedback) |  |  | <p>Feedback declares rule feedback entries.</p> |
| `assert` | [`Assert`](#json-kyverno-io-v1alpha1-Assert) | :white_check_mark: |  | <p>Assert is used to validate matching resources.</p> |

  