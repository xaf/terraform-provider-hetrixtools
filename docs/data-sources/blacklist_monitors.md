---
page_title: "hetrixtools_blacklist_monitors Data Source"
description: |-
  Lists HetrixTools blacklist monitors.
---

# hetrixtools_blacklist_monitors (Data Source)

Lists HetrixTools blacklist monitors.

## Example Usage

```terraform
data "hetrixtools_blacklist_monitors" "all" {
  per_page = "200"
}
```

## Schema

### Optional

- `cidr` (String) CIDR prefix for target filter.
- `exact_name` (String) Exact name filter.
- `exact_target` (String) Exact target filter.
- `listed` (String) Listed status.
- `name` (String) Name filter.
- `order` (String) Sort order.
- `order_by` (String) Sort field.
- `page` (String) Page number.
- `per_page` (String) Results per page.
- `target` (String) Target filter.
- `type` (String) Monitor type.

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
