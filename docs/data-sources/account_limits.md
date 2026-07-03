---
page_title: "hetrixtools_account_limits Data Source"
description: |-
  Reads HetrixTools account limits.
---

# hetrixtools_account_limits (Data Source)

Reads HetrixTools account limits.

## Example Usage

```terraform
data "hetrixtools_account_limits" "current" {}
```

## Schema

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
