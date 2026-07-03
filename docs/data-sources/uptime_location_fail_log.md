---
page_title: "hetrixtools_uptime_location_fail_log Data Source"
description: |-
  Gets a HetrixTools uptime monitor location fail log.
---

# hetrixtools_uptime_location_fail_log (Data Source)

Gets a HetrixTools uptime monitor location fail log.

## Example Usage

```terraform
data "hetrixtools_uptime_location_fail_log" "example" {
  monitor_id = hetrixtools_uptime_monitor.web.id
}
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.

### Optional

- `minutes` (String) Number of minutes with log entries to return.
- `timestamp` (String) Timestamp to start scanning from.

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
