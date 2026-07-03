---
page_title: "hetrixtools_server_agent Resource"
description: |-
  Attaches a HetrixTools server monitoring agent to an uptime monitor.
---

# hetrixtools_server_agent (Resource)

Attaches a HetrixTools server monitoring agent to an uptime monitor.

## Example Usage

```terraform
resource "hetrixtools_server_agent" "example" {
  monitor_id = hetrixtools_uptime_monitor.server.id
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_server_agent.example 00000000000000000000000000000000
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.

### Read-Only

- `agent_id` (String) HetrixTools server agent ID.
- `id` (String) Resource ID, equal to `monitor_id`.
