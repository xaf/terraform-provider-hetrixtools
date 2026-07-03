---
page_title: "hetrixtools_contact_lists Data Source"
description: |-
  Lists HetrixTools contact lists.
---

# hetrixtools_contact_lists (Data Source)

Lists HetrixTools contact lists.

## Example Usage

```terraform
data "hetrixtools_contact_lists" "all" {
  per_page = "200"
}
```

## Schema

### Optional

- `page` (String) Page number to return.
- `per_page` (String) Number of contact lists per page.

### Read-Only

- `json` (String) JSON response from the HetrixTools API.
