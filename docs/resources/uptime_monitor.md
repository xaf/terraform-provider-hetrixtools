---
page_title: "hetrixtools_uptime_monitor Resource"
description: |-
  Manages a HetrixTools uptime monitor.
---

# hetrixtools_uptime_monitor (Resource)

Manages a HetrixTools uptime monitor.

## Example Usage

```terraform
resource "hetrixtools_uptime_monitor" "web" {
  type                 = 1
  name                 = "Website"
  target               = "https://example.com"
  frequency            = 60
  timeout              = 10
  fails_before_alert   = 3
  failed_locations     = 2
  verify_ssl_host      = true
  verify_ssl_certificate = true

  locations = {
    ams = true
    nyc = true
  }
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_uptime_monitor.web 00000000000000000000000000000000
```

## Schema

### Required

- `name` (String) Monitor name.
- `type` (Number) Monitor type: `1` website, `2` ping/service, `3` SMTP, `9` server agent.

### Optional

- `alert_after` (String) Alert delay setting.
- `category` (String) Monitor category.
- `contact_list_id` (String) Contact list ID used for notifications.
- `cpu_public` (Boolean) Whether CPU details are public for server-agent monitors.
- `disk_public` (Boolean) Whether disk details are public for server-agent monitors.
- `extra_json` (String, Sensitive) Additional JSON fields merged into the uptime monitor payload for type-specific options like `Method`, `Keyword`, `HTTPCodes`, `SMTPUser`, or `SMTPPass`.
- `failed_locations` (Number) Number of failed locations required before alerting.
- `fails_before_alert` (Number) Number of failed checks required before alerting.
- `frequency` (Number) Check frequency.
- `grace` (Number) Grace period.
- `info_public` (Boolean) Whether server info is public for server-agent monitors.
- `locations` (Map of Boolean) Map of HetrixTools location code to enabled flag, e.g. `{ ams = true, nyc = false }`.
- `net_public` (Boolean) Whether network details are public for server-agent monitors.
- `public` (Boolean) Whether the monitor is public.
- `ram_public` (Boolean) Whether RAM details are public for server-agent monitors.
- `repeat_every` (String) Alert repeat interval.
- `repeat_times` (Number) Alert repeat count.
- `show_target` (Boolean) Whether to show the target publicly.
- `target` (String) Monitor target.
- `timeout` (Number) Check timeout.
- `verify_ssl_certificate` (Boolean) Whether to verify the SSL certificate.
- `verify_ssl_host` (Boolean) Whether to verify the SSL host.

### Read-Only

- `id` (String) Uptime monitor ID.
- `server_id` (String, Sensitive) Server agent ID returned for server-agent monitors.
