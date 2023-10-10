---
title: kyverno-json (v1alpha1)
content_type: tool-reference
package: json.kyverno.io/v1alpha1
auto_generated: true
---
<p>Package v1alpha1 contains API Schema definitions for the policy v1alpha1 API group</p>


## Resource Types 


- [Policy](#json-kyverno-io-v1alpha1-Policy)
  

## `Policy`     {#json-kyverno-io-v1alpha1-Policy}
    



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
<tr><td><code>apiVersion</code><br/>string</td><td><code>json.kyverno.io/v1alpha1</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>Policy</code></td></tr>
    
  
<tr><td><code>TypeMeta</code> <B>[Required]</B><br/>
<code>k8s.io/apimachinery/pkg/apis/meta/v1.TypeMeta</code>
</td>
<td>(Members of <code>TypeMeta</code> are embedded into this type.)
   <span class="text-muted">No description provided.</span></td>
</tr>
<tr><td><code>metadata</code> <B>[Required]</B><br/>
<code>k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta</code>
</td>
<td>
   <span class="text-muted">No description provided.</span>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr><td><code>spec</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-PolicySpec"><code>PolicySpec</code></a>
</td>
<td>
   <span class="text-muted">No description provided.</span></td>
</tr>
</tbody>
</table>

## `Any`     {#json-kyverno-io-v1alpha1-Any}
    

**Appears in:**

- [ContextEntry](#json-kyverno-io-v1alpha1-ContextEntry)



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>Value</code> <B>[Required]</B><br/>
<code>interface{}</code>
</td>
<td>(Members of <code>Value</code> are embedded into this type.)
   <span class="text-muted">No description provided.</span></td>
</tr>
</tbody>
</table>

## `Assertions`     {#json-kyverno-io-v1alpha1-Assertions}
    
(Alias of `[]github.com/kyverno/kyverno-json/pkg/apis/v1alpha1.Any`)

**Appears in:**

- [Match](#json-kyverno-io-v1alpha1-Match)





## `ContextEntry`     {#json-kyverno-io-v1alpha1-ContextEntry}
    

**Appears in:**

- [Rule](#json-kyverno-io-v1alpha1-Rule)


<p>ContextEntry adds variables and data sources to a rule Context.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Name is the variable name.</p>
</td>
</tr>
<tr><td><code>variable</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Any"><code>Any</code></a>
</td>
<td>
   <p>Variable defines an arbitrary JMESPath context variable that can be defined inline.</p>
</td>
</tr>
</tbody>
</table>

## `Match`     {#json-kyverno-io-v1alpha1-Match}
    

**Appears in:**

- [Rule](#json-kyverno-io-v1alpha1-Rule)

- [Validation](#json-kyverno-io-v1alpha1-Validation)



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>any</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Assertions"><code>Assertions</code></a>
</td>
<td>
   <p>Any allows specifying resources which will be ORed.</p>
</td>
</tr>
<tr><td><code>all</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Assertions"><code>Assertions</code></a>
</td>
<td>
   <p>All allows specifying resources which will be ANDed.</p>
</td>
</tr>
</tbody>
</table>

## `PolicySpec`     {#json-kyverno-io-v1alpha1-PolicySpec}
    

**Appears in:**

- [Policy](#json-kyverno-io-v1alpha1-Policy)



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>rules</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Rule"><code>[]Rule</code></a>
</td>
<td>
   <p>Rules is a list of Rule instances. A Policy contains multiple rules and each rule can validate, mutate, or generate resources.</p>
</td>
</tr>
</tbody>
</table>

## `Rule`     {#json-kyverno-io-v1alpha1-Rule}
    

**Appears in:**

- [PolicySpec](#json-kyverno-io-v1alpha1-PolicySpec)



<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>name</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Name is a label to identify the rule, It must be unique within the policy.</p>
</td>
</tr>
<tr><td><code>context</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-ContextEntry"><code>[]ContextEntry</code></a>
</td>
<td>
   <p>Context defines variables and data sources that can be used during rule execution.</p>
</td>
</tr>
<tr><td><code>match</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Match"><code>Match</code></a>
</td>
<td>
   <p>Match defines when this policy rule should be applied. The match
criteria can include resource information (e.g. kind, name, namespace, labels)
and admission review request information like the user name or role.
At least one kind is required.</p>
</td>
</tr>
<tr><td><code>exclude</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Match"><code>Match</code></a>
</td>
<td>
   <p>Exclude defines when this policy rule should not be applied. The exclude
criteria can include resource information (e.g. kind, name, namespace, labels)
and admission review request information like the name or role.</p>
</td>
</tr>
<tr><td><code>validate</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Validation"><code>Validation</code></a>
</td>
<td>
   <p>Validation is used to validate matching resources.</p>
</td>
</tr>
</tbody>
</table>

## `Validation`     {#json-kyverno-io-v1alpha1-Validation}
    

**Appears in:**

- [Rule](#json-kyverno-io-v1alpha1-Rule)


<p>Validation defines checks to be performed on matching resources.</p>


<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    
  
<tr><td><code>message</code> <B>[Required]</B><br/>
<code>string</code>
</td>
<td>
   <p>Message specifies a custom message to be displayed on failure.</p>
</td>
</tr>
<tr><td><code>assert</code> <B>[Required]</B><br/>
<a href="#json-kyverno-io-v1alpha1-Match"><code>Match</code></a>
</td>
<td>
   <p>Assert specifies an overlay-style pattern used to check resources.</p>
</td>
</tr>
</tbody>
</table>
  