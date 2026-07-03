---
page_title: "hetrixtools_blacklist_monitor Resource"
description: |-
  Manages a HetrixTools blacklist monitor.
---

# hetrixtools_blacklist_monitor (Resource)

Manages a HetrixTools blacklist monitor.

## Example Usage

```terraform
resource "hetrixtools_blacklist_monitor" "example" {
  target          = "203.0.113.10"
  label           = "Example IP"
  contact_list_id = "00000000000000000000000000000000"
}
```

## Import

Import by target value:

```shell
terraform import hetrixtools_blacklist_monitor.example 203.0.113.10
```

## Schema

### Required

- `target` (String) IP address, IP range, CIDR block, or domain name.

### Optional

- `contact_list_id` (String) Contact list ID used for notifications.
- `label` (String) Monitor label.

### Read-Only

- `id` (String) Blacklist monitor ID.
