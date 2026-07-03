---
page_title: "hetrixtools_uptime_server_agent Data Source"
description: |-
  Gets the HetrixTools server agent attached to an uptime monitor.
---

# hetrixtools_uptime_server_agent (Data Source)

Gets the HetrixTools server agent attached to an uptime monitor.

## Example Usage

```terraform
data "hetrixtools_uptime_server_agent" "example" {
  monitor_id = hetrixtools_uptime_monitor_heartbeat.server.id
}
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
