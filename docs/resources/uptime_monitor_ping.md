---
page_title: "hetrixtools_uptime_monitor_ping Resource"
description: |-
  Manages a HetrixTools ping uptime monitor.
---

# hetrixtools_uptime_monitor_ping (Resource)

Manages a HetrixTools ping uptime monitor.

Create, update, delete, and import operations are paced by the shared HetrixTools client to avoid API rate limits. See the client rate-limiting notes and HetrixTools API references in the provider index.

## Example Usage

```terraform
resource "hetrixtools_uptime_monitor_ping" "gateway" {
  name   = "Gateway"
  target = "198.51.100.10"

  locations = ["amsterdam", "new_york"]
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_uptime_monitor_ping.gateway 00000000000000000000000000000000
```

Imported resources should still be made self-sufficient in Terraform configuration so the monitor can be recreated from scratch.

## Defaults

Most optional arguments are also computed. If omitted, the provider does not send that field during create/update and keeps the value returned by HetrixTools on read/import. For new monitors, HetrixTools applies its own defaults.

## Schema

### Required

- `name` (String) Monitor name.
- `target` (String) Hostname or IP address to ping.

### Optional

- `alert_after` (String) Delay before sending an alert, such as `5m`.
- `category` (String) Monitor category.
- `contact_list_id` (String) Contact list ID used for notifications.
- `failed_locations` (Number) Number of failed monitoring locations required before alerting.
- `fails_before_alert` (Number) Number of failed checks required before alerting.
- `frequency` (Number) Check frequency.
- `locations` (Set of String) Canonical HetrixTools location names enabled for this monitor. Supported values are `new_york`, `san_francisco`, `dallas`, `amsterdam`, `london`, `frankfurt`, `singapore`, `sydney`, `sao_paulo`, `tokyo`, `mumbai`, and `warsaw`.
- `public` (Boolean) Whether the monitor has a public report.
- `repeat_every` (String) Alert repeat interval, such as `60m`.
- `repeat_times` (Number) Number of times to repeat alerts.
- `show_target` (Boolean) Whether to show the monitor target publicly.
- `timeout` (Number) Check timeout.

### Read-Only

- `id` (String) Uptime monitor ID.
