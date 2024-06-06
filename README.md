# Bucket Lifecycle Function

## Required Dynamic Group

```text
ALL {resource.type = 'fnfunc', resource.compartment.id = 'ocid1.compartment.oc1..xyz'}
```

## Required Policies

```text
Allow dynamic-group <function-dynamic-group> to use buckets in tenancy
Allow dynamic-group <function-dynamic-group> to manage objects in tenancy
```
