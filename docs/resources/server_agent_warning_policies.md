---
page_title: "hetrixtools_server_agent_warning_policies Resource"
description: |-
  Manages server-agent warning policies for a HetrixTools uptime monitor.
---

# hetrixtools_server_agent_warning_policies (Resource)

Manages server-agent warning policies for an uptime monitor. The policy payload mirrors the HetrixTools API JSON schema.

## Example Usage

```terraform
resource "hetrixtools_server_agent_warning_policies" "example" {
  monitor_id = hetrixtools_uptime_monitor.server.id

  policies_json = jsonencode({
    cpu = {
      enabled = true
      warning = 80
    }
  })
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_server_agent_warning_policies.example 00000000000000000000000000000000
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.
- `policies_json` (String) JSON object sent to HetrixTools to replace the monitor's server-agent warning policies.

### Read-Only

- `id` (String) Resource ID, equal to `monitor_id`.

## Delete Behavior

HetrixTools does not expose a delete or reset endpoint for server-agent warning policies. Destroying this resource removes it from Terraform state but leaves the remote policies in place.
