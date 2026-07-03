---
page_title: "hetrixtools_uptime_downtimes Data Source"
description: |-
  Lists HetrixTools uptime monitor downtimes.
---

# hetrixtools_uptime_downtimes (Data Source)

Lists HetrixTools uptime monitor downtimes.

## Example Usage

```terraform
data "hetrixtools_uptime_downtimes" "example" {
  monitor_id = hetrixtools_uptime_monitor_http.web.id
}
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.

### Optional

- `page` (String) Page number.
- `per_page` (String) Results per page.
- `start_after` (String) Only downtimes started after this timestamp.
- `start_before` (String) Only downtimes started before this timestamp.

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
