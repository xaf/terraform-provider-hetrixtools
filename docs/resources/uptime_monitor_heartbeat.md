---
page_title: "hetrixtools_uptime_monitor_heartbeat Resource"
description: |-
  Manages a HetrixTools heartbeat uptime monitor.
---

# hetrixtools_uptime_monitor_heartbeat (Resource)

Manages a HetrixTools heartbeat uptime monitor, also called a server-agent monitor in parts of the HetrixTools API.

Create, update, delete, and import operations are paced by the shared HetrixTools client to avoid API rate limits. See the client rate-limiting notes and HetrixTools API references in the provider index.

## Example Usage

```terraform
resource "hetrixtools_uptime_monitor_heartbeat" "server" {
  name  = "Server Agent"
  grace = 120
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_uptime_monitor_heartbeat.server 00000000000000000000000000000000
```

Imported resources should still be made self-sufficient in Terraform configuration so the monitor can be recreated from scratch.

## Defaults

Most optional arguments are also computed. If omitted, the provider does not send that field during create/update and keeps the value returned by HetrixTools on read/import. For new monitors, HetrixTools applies its own defaults.

## Schema

### Required

- `name` (String) Monitor name.

### Optional

- `alert_after` (String) Delay before sending an alert, such as `5m`.
- `category` (String) Monitor category.
- `contact_list_id` (String) Contact list ID used for notifications.
- `cpu_public` (Boolean) Whether CPU details are public.
- `disk_public` (Boolean) Whether disk details are public.
- `fails_before_alert` (Number) Number of failed checks required before alerting.
- `frequency` (Number) Check frequency.
- `grace` (Number) Grace period for heartbeat monitors.
- `info_public` (Boolean) Whether server info is public.
- `net_public` (Boolean) Whether network details are public.
- `public` (Boolean) Whether the monitor has a public report.
- `ram_public` (Boolean) Whether RAM details are public.
- `repeat_every` (String) Alert repeat interval, such as `60m`.
- `repeat_times` (Number) Number of times to repeat alerts.
- `show_target` (Boolean) Whether to show the monitor target publicly.
- `timeout` (Number) Check timeout.

### Read-Only

- `id` (String) Uptime monitor ID.
- `server_id` (String) Server agent ID returned for heartbeat monitors.
