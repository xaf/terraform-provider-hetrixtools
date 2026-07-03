---
page_title: "hetrixtools_scheduled_maintenance Resource"
description: |-
  Manages a HetrixTools scheduled maintenance window.
---

# hetrixtools_scheduled_maintenance (Resource)

Manages a HetrixTools scheduled maintenance window for an uptime monitor.

Create, delete, and import operations are paced by the shared HetrixTools client to avoid API rate limits. See the client rate-limiting notes and HetrixTools API references in the provider index.

## Example Usage

```terraform
resource "hetrixtools_scheduled_maintenance" "example" {
  monitor_id         = "00000000000000000000000000000000"
  start              = "2026-07-02 10:00"
  end                = "2026-07-02 11:00"
  timezone           = "UTC"
  with_notifications = true
}
```

## Import

Import by scheduled maintenance ID:

```shell
terraform import hetrixtools_scheduled_maintenance.example 00000000000000000000000000000000
```

## Schema

### Required

- `end` (String) Maintenance end time as `YYYY-MM-DD HH:MM`.
- `monitor_id` (String) Uptime monitor ID.
- `start` (String) Maintenance start time as `YYYY-MM-DD HH:MM`.
- `timezone` (String) Maintenance timezone.

### Optional

- `recurring_time` (Number) Recurring interval value.
- `recurring_time_type` (String) One of `hour`, `day`, `week`, `month`, or `year` when `recurring_time` is set.
- `with_notifications` (Boolean) Whether HetrixTools should send maintenance notifications.

### Read-Only

- `id` (String) Scheduled maintenance ID.
- `recurring` (Boolean) Whether the maintenance window is recurring.
