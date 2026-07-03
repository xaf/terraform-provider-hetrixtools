---
page_title: "hetrixtools_uptime_monitors Data Source"
description: |-
  Lists HetrixTools uptime monitors.
---

# hetrixtools_uptime_monitors (Data Source)

Lists HetrixTools uptime monitors.

## Example Usage

```terraform
data "hetrixtools_uptime_monitors" "all" {
  per_page = "200"
}
```

## Schema

### Optional

- `category` (String) Category filter.
- `id` (String) Monitor ID.
- `monitor_status` (String) Monitor status.
- `name` (String) Name filter.
- `order` (String) Sort order.
- `order_by` (String) Sort field.
- `page` (String) Page number.
- `per_page` (String) Results per page.
- `target` (String) Target filter.
- `type` (String) Monitor type.
- `uptime_status` (String) Uptime status.

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
