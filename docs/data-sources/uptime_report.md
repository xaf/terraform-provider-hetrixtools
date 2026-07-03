---
page_title: "hetrixtools_uptime_report Data Source"
description: |-
  Gets a HetrixTools uptime monitor report.
---

# hetrixtools_uptime_report (Data Source)

Gets a HetrixTools uptime monitor report.

## Example Usage

```terraform
data "hetrixtools_uptime_report" "example" {
  monitor_id = hetrixtools_uptime_monitor_http.web.id
  days       = "7"
}
```

## Schema

### Required

- `monitor_id` (String) Uptime monitor ID.

### Optional

- `days` (String) Number of recent days to display.
- `hourly_stats` (String) Whether to include hourly stats.
- `month` (String) Month in `YYYY-MM` format.
- `timezone` (String) Report timezone.

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
