---
page_title: "hetrixtools_uptime_monitor_http Resource"
description: |-
  Manages a HetrixTools HTTP uptime monitor.
---

# hetrixtools_uptime_monitor_http (Resource)

Manages a HetrixTools HTTP uptime monitor.

## Example Usage

```terraform
resource "hetrixtools_uptime_monitor_http" "web" {
  name   = "Website"
  target = "https://example.com"

  locations           = ["amsterdam", "new_york"]
  http_method         = "GET"
  keyword             = "healthy"
  accepted_http_codes = [200, 204]
}
```

## Import

Import by uptime monitor ID:

```shell
terraform import hetrixtools_uptime_monitor_http.web 00000000000000000000000000000000
```

Imported resources should still be made self-sufficient in Terraform configuration so the monitor can be recreated from scratch.

## Defaults

Most optional arguments are also computed. If omitted, the provider does not send that field during create/update and keeps the value returned by HetrixTools on read/import. For new monitors, HetrixTools applies its own defaults.

## Schema

### Required

- `name` (String) Monitor name.
- `target` (String) URL to check, for example `https://example.com`.

### Optional

- `accepted_http_codes` (List of Number) HTTP status codes accepted as healthy.
- `alert_after` (String) Delay before sending an alert, such as `5m`.
- `category` (String) Monitor category.
- `contact_list_id` (String) Contact list ID used for notifications.
- `failed_locations` (Number) Number of failed monitoring locations required before alerting.
- `fails_before_alert` (Number) Number of failed checks required before alerting.
- `frequency` (Number) Check frequency.
- `http_method` (String) HTTP request method, for example `GET`.
- `keyword` (String) Keyword expected in the response body.
- `locations` (Set of String) Canonical HetrixTools location names enabled for this monitor. Supported values are `new_york`, `san_francisco`, `dallas`, `amsterdam`, `london`, `frankfurt`, `singapore`, `sydney`, `sao_paulo`, `tokyo`, `mumbai`, and `warsaw`.
- `max_redirects` (Number) Maximum HTTP redirects to follow.
- `public` (Boolean) Whether the monitor has a public report.
- `repeat_every` (String) Alert repeat interval, such as `60m`.
- `repeat_times` (Number) Number of times to repeat alerts.
- `show_target` (Boolean) Whether to show the monitor target publicly.
- `timeout` (Number) Check timeout.
- `verify_ssl_certificate` (Boolean) Whether to verify the SSL certificate.
- `verify_ssl_host` (Boolean) Whether to verify the SSL hostname.

### Read-Only

- `id` (String) Uptime monitor ID.
