---
page_title: "hetrixtools_account_limits Data Source"
description: |-
  Reads HetrixTools account limits.
---

# hetrixtools_account_limits (Data Source)

Reads HetrixTools account limits.

The `json` attribute contains the HetrixTools API response. Use Terraform's `jsondecode()` to access individual fields.

## Example Usage

```terraform
data "hetrixtools_account_limits" "current" {}
```

## Schema

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
