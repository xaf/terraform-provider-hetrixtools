---
page_title: "hetrixtools_uptime_monitor_smtp Resource"
description: |-
  Manages a HetrixTools SMTP uptime monitor.
---

# hetrixtools_uptime_monitor_smtp (Resource)

Manages a HetrixTools SMTP uptime monitor.

Create, update, delete, and import operations are paced by the shared HetrixTools client to avoid API rate limits. See the client rate-limiting notes and HetrixTools API references in the provider index.

## Example Usage

```terraform
resource "hetrixtools_uptime_monitor_smtp" "mail" {
  name   = "SMTP"
  target = "smtp.example.com"
  port   = 587

  locations = ["amsterdam", "new_york"]
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_uptime_monitor_smtp.mail 00000000000000000000000000000000
```

Imported resources should still be made self-sufficient in Terraform configuration. Imported SMTP monitors must declare `port` in HCL so the resource can be recreated from scratch.

## Defaults

Most optional arguments are also computed. If omitted, the provider does not send that field during create/update and keeps the value returned by HetrixTools on read/import. For new monitors, HetrixTools applies its own defaults.

## Schema

### Required

- `name` (String) Monitor name.
- `port` (Number) SMTP port.
- `target` (String) SMTP hostname.

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
- `smtp_password` (String, Sensitive) SMTP password. Must be set together with `smtp_user`. HetrixTools does not return this value on read, so Terraform preserves the configured state value.
- `smtp_user` (String) SMTP username. Must be set together with `smtp_password`.
- `timeout` (Number) Check timeout.
- `verify_ssl_certificate` (Boolean) Whether to verify the SSL certificate.
- `verify_ssl_host` (Boolean) Whether to verify the SSL hostname.

### Read-Only

- `id` (String) Uptime monitor ID.
