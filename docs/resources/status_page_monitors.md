---
page_title: "hetrixtools_status_page_monitors Resource"
description: |-
  Manages the exact monitor set on a HetrixTools status page.
---

# hetrixtools_status_page_monitors (Resource)

Manages the exact monitor set on a HetrixTools status page.

## Example Usage

```terraform
resource "hetrixtools_status_page_monitors" "example" {
  status_page_id = "00000000000000000000000000000000"

  monitor_ids = [
    hetrixtools_uptime_monitor_http.web.id,
  ]
}
```

## Import

Import by status page ID:

```shell
terraform import hetrixtools_status_page_monitors.example 00000000000000000000000000000000
```

## Schema

### Required

- `monitor_ids` (Set of String) Exact set of uptime monitor IDs assigned to the status page.
- `status_page_id` (String) Status page ID.

### Read-Only

- `id` (String) Resource ID, equal to `status_page_id`.
