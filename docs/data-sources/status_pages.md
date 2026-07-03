---
page_title: "hetrixtools_status_pages Data Source"
description: |-
  Lists HetrixTools status pages.
---

# hetrixtools_status_pages (Data Source)

Lists HetrixTools status pages.

## Example Usage

```terraform
data "hetrixtools_status_pages" "all" {
  per_page = "200"
}
```

## Schema

### Optional

- `name` (String) Status page name filter.
- `page` (String) Page number.
- `per_page` (String) Results per page.

### Read-Only

- `json` (String) JSON response from the HetrixTools API. Use Terraform's `jsondecode()` to access individual fields.
