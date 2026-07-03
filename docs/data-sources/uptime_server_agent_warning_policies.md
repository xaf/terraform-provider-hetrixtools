---
page_title: "hetrixtools_uptime_server_agent_warning_policies Data Source"
description: |-
  Gets HetrixTools server agent warning policies for an uptime monitor.
---

# hetrixtools_uptime_server_agent_warning_policies (Data Source)

Gets HetrixTools server agent warning policies for an uptime monitor.

## Example Usage

```terraform
data "hetrixtools_uptime_server_agent_warning_policies" "example" {
  monitor_id = hetrixtools_uptime_monitor_heartbeat.server.id
}
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
