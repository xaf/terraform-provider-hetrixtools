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
  type   = "http"
  name   = "Website"
  target = "https://example.com"

  locations           = ["amsterdam", "new_york"]
  http_method         = "GET"
  keyword             = "healthy"
  accepted_http_codes = [200, 204]
}

resource "hetrixtools_uptime_monitor" "ping" {
  type   = "ping"
  name   = "Gateway"
  target = "198.51.100.10"
}

resource "hetrixtools_uptime_monitor" "smtp" {
  type   = "smtp"
  name   = "SMTP"
  target = "smtp.example.com"
  port   = 587
}

resource "hetrixtools_uptime_monitor" "heartbeat" {
  type  = "heartbeat"
  name  = "Server Agent"
  grace = 120
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_uptime_monitor.web 00000000000000000000000000000000
```

Imported resources should still be made self-sufficient in Terraform configuration. For example, imported SMTP monitors must declare `port` in HCL so the resource can be recreated from scratch.

## Defaults

Most optional arguments are also computed. If omitted, the provider does not send that field during create/update and keeps the value returned by HetrixTools on read/import. For new monitors, HetrixTools applies its own defaults.

## Schema

### Required for All Monitor Types

- `name` (String) Monitor name.
- `type` (String) Monitor type: `http`, `ping`, `smtp`, or `heartbeat`. The provider maps these values to HetrixTools API type IDs internally. Changing this value forces replacement.

### Minimum Fields by Monitor Type

All monitors require `type` and `name`.

| Type | Additional required fields | Notes |
|---|---|---|
| `http` | `target` | URL to check, for example `https://example.com`. |
| `ping` | `target` | Hostname or IP address to ping. |
| `smtp` | `target`, `port` | `port` is validated as required for SMTP monitors. |
| `heartbeat` | None | Do not set `target`, `locations`, or `failed_locations`. `server_id` is returned after creation. |

### Optional Common Fields

These fields are valid for all monitor types unless HetrixTools rejects a value for a specific account or plan:

- `alert_after` (String) Delay before sending an alert, such as `5m`.
- `category` (String) Monitor category.
- `contact_list_id` (String) Contact list ID used for notifications.
- `fails_before_alert` (Number) Number of failed checks required before alerting.
- `frequency` (Number) Check frequency.
- `public` (Boolean) Whether the monitor has a public report.
- `repeat_every` (String) Alert repeat interval, such as `60m`.
- `repeat_times` (Number) Number of times to repeat alerts.
- `show_target` (Boolean) Whether to show the monitor target publicly.
- `timeout` (Number) Check timeout.

### Optional Fields for Non-Heartbeat Monitors

Valid for `http`, `ping`, and `smtp`; invalid for `heartbeat`:

- `failed_locations` (Number) Number of failed monitoring locations required before alerting.
- `locations` (Set of String) Canonical HetrixTools location names enabled for this monitor. Supported values are `new_york`, `san_francisco`, `dallas`, `amsterdam`, `london`, `frankfurt`, `singapore`, `sydney`, `sao_paulo`, `tokyo`, `mumbai`, and `warsaw`. The provider converts these to legacy v2 API location codes when sending requests.

### HTTP-Only Optional Fields

Only valid when `type = "http"`:

- `accepted_http_codes` (List of Number) HTTP status codes accepted as healthy.
- `http_method` (String) HTTP request method, for example `GET`.
- `keyword` (String) Keyword expected in the response body.
- `max_redirects` (Number) Maximum HTTP redirects to follow.

### SMTP-Only Fields

Only valid when `type = "smtp"`:

- `port` (Number) SMTP port. Required.
- `smtp_user` (String) SMTP username. Must be set together with `smtp_password`.
- `smtp_password` (String, Sensitive) SMTP password. Must be set together with `smtp_user`. HetrixTools does not return this value on read, so Terraform preserves the configured state value.

### HTTP and SMTP Optional Fields

Only valid when `type` is `http` or `smtp`:

- `verify_ssl_certificate` (Boolean) Whether to verify the SSL certificate.
- `verify_ssl_host` (Boolean) Whether to verify the SSL hostname.

### Heartbeat-Only Optional Fields

Only valid when `type = "heartbeat"`:

- `grace` (Number) Grace period for heartbeat monitors.
- `info_public` (Boolean) Whether server info is public.
- `cpu_public` (Boolean) Whether CPU details are public.
- `ram_public` (Boolean) Whether RAM details are public.
- `disk_public` (Boolean) Whether disk details are public.
- `net_public` (Boolean) Whether network details are public.

### Validation Notes

- `type` must be one of `http`, `ping`, `smtp`, or `heartbeat`.
- `port` is required for `smtp` monitors and invalid for all other types.
- `smtp_user` and `smtp_password` must either both be set or both be omitted.
- `target` is required for `http`, `ping`, and `smtp` monitors.
- `target`, `locations`, and `failed_locations` are not supported for `heartbeat` monitors.
- `http_method`, `max_redirects`, `keyword`, and `accepted_http_codes` are valid only for `http` monitors.
- `verify_ssl_certificate` and `verify_ssl_host` are valid only for `http` and `smtp` monitors.
- Heartbeat public-detail fields and `grace` are valid only for `heartbeat` monitors.
- `locations` must use canonical HetrixTools location names such as `amsterdam` or `new_york`; legacy API location codes are not accepted in Terraform configuration.
- Setting a type-specific boolean or number on an unsupported type is invalid even if the value is `false` or `0`.

### Read-Only Attributes

- `id` (String) Uptime monitor ID.
- `server_id` (String, Sensitive) Server agent ID returned for heartbeat monitors.
