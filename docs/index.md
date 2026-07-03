---
page_title: "HetrixTools Provider"
description: |-
  Terraform provider for managing HetrixTools uptime monitors, blacklist monitors, scheduled maintenance, status page monitor membership, and server-agent settings.
---

# HetrixTools Provider

The HetrixTools provider manages resources exposed by the HetrixTools API.

This project is not affiliated with, endorsed by, or supported by HetrixTools. It is provided without any guarantee that the HetrixTools API or this provider will behave as expected for your account.

Canonical HCL documentation is published on the Terraform Registry: https://registry.terraform.io/providers/xaf/hetrixtools/latest/docs.

Terraform users should use the Registry documentation on this page. The reusable Go client used internally by the provider is documented on pkg.go.dev: https://pkg.go.dev/github.com/xaf/terraform-provider-hetrixtools/client.

## Example Usage

```terraform
terraform {
  required_providers {
    hetrixtools = {
      source  = "xaf/hetrixtools"
      version = "~> 0.1.6"
    }
  }
}

provider "hetrixtools" {
  api_token = var.hetrixtools_api_token
}
```

## Authentication

The provider can be configured with an API token argument or an environment variable.

```terraform
provider "hetrixtools" {
  api_token = var.hetrixtools_api_token
}
```

Environment variables:

- `HETRIXTOOLS_API_TOKEN`
- `HETRIXTOOLS_BASE_URL` for the API root URL. The provider appends `/v2` and `/v3`.
- `HETRIXTOOLS_BASE_URL_V2` to override only the v2 API base URL.
- `HETRIXTOOLS_BASE_URL_V3` to override only the v3 API base URL.

## Client Rate Limiting

The provider delegates request pacing, retries, and rate-limit reset handling to the reusable Go client. This behavior is based on the HetrixTools API overview and v3 reference:

- https://docs.hetrixtools.com/understanding-our-apis/
- https://docs.hetrixtools.com/api-v3/

Legacy v2 token-path endpoints are paced through one shared limiter. v3 REST endpoints are paced through both a user-level limiter and a normalized method/path endpoint limiter.

Client package documentation: https://pkg.go.dev/github.com/xaf/terraform-provider-hetrixtools/client.

This protects HetrixTools API limits, but large applies may take longer when many resources are created, updated, imported, or read. Users should not need to lower Terraform parallelism just to avoid HetrixTools rate limits.

## Schema

### Optional

- `api_token` (String, Sensitive) HetrixTools API bearer token. Can also be set with `HETRIXTOOLS_API_TOKEN`.
- `base_url` (String) HetrixTools API root URL. The provider appends `/v2` and `/v3`. Can also be set with `HETRIXTOOLS_BASE_URL`.
- `base_url_v2` (String) HetrixTools v2 API base URL. Overrides `base_url` for token-path endpoints. Can also be set with `HETRIXTOOLS_BASE_URL_V2`.
- `base_url_v3` (String) HetrixTools v3 API base URL. Overrides `base_url` for REST endpoints. Can also be set with `HETRIXTOOLS_BASE_URL_V3`.

## Resources

- `hetrixtools_blacklist_monitor`
- `hetrixtools_scheduled_maintenance`
- `hetrixtools_server_agent`
- `hetrixtools_server_agent_warning_policies`
- `hetrixtools_status_page_monitors`
- `hetrixtools_uptime_monitor_heartbeat`
- `hetrixtools_uptime_monitor_http`
- `hetrixtools_uptime_monitor_ping`
- `hetrixtools_uptime_monitor_smtp`

## Data Sources

- `hetrixtools_account_limits`
- `hetrixtools_blacklist_monitors`
- `hetrixtools_blacklist_report`
- `hetrixtools_blacklists`
- `hetrixtools_contact_lists`
- `hetrixtools_scheduled_maintenances`
- `hetrixtools_status_pages`
- `hetrixtools_uptime_downtimes`
- `hetrixtools_uptime_location_fail_log`
- `hetrixtools_uptime_monitors`
- `hetrixtools_uptime_report`
- `hetrixtools_uptime_server_agent`
- `hetrixtools_uptime_server_agent_warning_policies`
