---
page_title: "hetrixtools_blacklists Data Source"
description: |-
  Lists HetrixTools blacklist providers.
---

# hetrixtools_blacklists (Data Source)

Lists HetrixTools blacklist providers.

## Example Usage

```terraform
data "hetrixtools_blacklists" "all" {}
```

## Schema

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
