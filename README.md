# Go HetrixTools

This repository contains a reusable Go client for HetrixTools plus a Terraform Plugin Framework provider built on top of that client.

This project is not affiliated with, endorsed by, or supported by HetrixTools. It is provided without any guarantee that the HetrixTools API or this provider will behave as expected for your account.

The client lives in `client/` and exposes semantic methods such as `GetUptimeMonitor`, `CreateUptimeMonitor`, `CreateBlacklistMonitor`, and `CreateScheduledMaintenance`. HetrixTools API-version selection, bearer auth, and token-in-path behavior stay inside the client.

The Terraform provider lives in `terraform-provider/` and calls the semantic client methods instead of hard-coding HetrixTools API versions.

## Supported Surface

Resources:

- `hetrixtools_blacklist_monitor`
- `hetrixtools_uptime_monitor_http`
- `hetrixtools_uptime_monitor_ping`
- `hetrixtools_uptime_monitor_smtp`
- `hetrixtools_uptime_monitor_heartbeat`
- `hetrixtools_scheduled_maintenance`
- `hetrixtools_status_page_monitors`
- `hetrixtools_server_agent`
- `hetrixtools_server_agent_warning_policies`

Typed read-only data sources:

- `hetrixtools_account_limits`
- `hetrixtools_blacklists`
- `hetrixtools_contact_lists`
- `hetrixtools_blacklist_monitors`
- `hetrixtools_blacklist_report`
- `hetrixtools_uptime_monitors`
- `hetrixtools_uptime_report`
- `hetrixtools_uptime_downtimes`
- `hetrixtools_uptime_location_fail_log`
- `hetrixtools_uptime_server_agent`
- `hetrixtools_uptime_server_agent_warning_policies`
- `hetrixtools_status_pages`
- `hetrixtools_scheduled_maintenances`

The provider combines the latest REST endpoints with older HetrixTools APIs where needed so uptime and blacklist monitor management are available as Terraform resources.

## Terraform Documentation

HCL documentation is available on the Terraform Registry:

- https://registry.terraform.io/providers/xaf/hetrixtools/latest/docs

The same Terraform Registry-format documentation is maintained in `docs/`:

- Provider configuration: `docs/index.md`
- Resources: `docs/resources/`
- Data sources: `docs/data-sources/`

The resource docs are written in Terraform Registry format and are the canonical HCL reference for this provider.

## Go Client Documentation

The Go client package is documented in `client/`, and the published package documentation is available at https://pkg.go.dev/github.com/xaf/terraform-provider-hetrixtools/client.

To preview the package docs locally:

```bash
go install golang.org/x/pkgsite/cmd/pkgsite@latest
pkgsite -http=:6060
```

Then open `http://localhost:6060/github.com/xaf/terraform-provider-hetrixtools/client`.

## Configuration

```hcl
provider "hetrixtools" {
  api_token = var.hetrixtools_api_token
}
```

You can also use environment variables:

- `HETRIXTOOLS_API_TOKEN`
- `HETRIXTOOLS_BASE_URL` for the API root URL. The provider appends `/v2` and `/v3`.
- `HETRIXTOOLS_BASE_URL_V2` to override only the v2 API base URL.
- `HETRIXTOOLS_BASE_URL_V3` to override only the v3 API base URL.

## Example

```hcl
data "hetrixtools_uptime_monitors" "web" {
  name = "example"
}

resource "hetrixtools_scheduled_maintenance" "window" {
  monitor_id         = "00000000000000000000000000000000"
  start              = "2026-07-02 10:00"
  end                = "2026-07-02 11:00"
  timezone           = "UTC"
  with_notifications = true
}
```

## Development

```bash
go mod tidy
CGO_ENABLED=0 go test ./...
CGO_ENABLED=0 go build ./...
```

## Releases

Terraform provider binaries are built with GoReleaser. Push a `vX.Y.Z` tag to create a GitHub release with zipped provider binaries named like:

```text
terraform-provider-hetrixtools_X.Y.Z_linux_amd64.zip
```

These assets are consumed by the Terraform infrastructure workspace cache script.

## License

Apache-2.0. See `LICENSE`.
