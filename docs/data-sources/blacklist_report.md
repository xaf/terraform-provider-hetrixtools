---
page_title: "hetrixtools_blacklist_report Data Source"
description: |-
  Gets a HetrixTools blacklist monitor report.
---

# hetrixtools_blacklist_report (Data Source)

Gets a HetrixTools blacklist monitor report.

## Example Usage

```terraform
data "hetrixtools_blacklist_report" "example" {
  identifier = "203.0.113.10"
}
```

## Schema

### Required

- `identifier` (String) Blacklist monitor ID, IP address, or hostname.

### Optional

- `date` (String) Report date in `YYYY-MM-DD` format.

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
