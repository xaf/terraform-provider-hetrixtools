---
page_title: "hetrixtools_scheduled_maintenances Data Source"
description: |-
  Lists HetrixTools scheduled maintenances.
---

# hetrixtools_scheduled_maintenances (Data Source)

Lists HetrixTools scheduled maintenances.

## Example Usage

```terraform
data "hetrixtools_scheduled_maintenances" "all" {
  monitor_id = hetrixtools_uptime_monitor_http.web.id
}
```

## Schema

### Optional

- `monitor_id` (String) Optional monitor ID filter.
- `page` (String) Page number.
- `per_page` (String) Results per page.

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
