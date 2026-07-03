---
page_title: "HetrixTools Provider"
description: |-
  Terraform provider for managing HetrixTools uptime monitors, blacklist monitors, scheduled maintenance, status page monitor membership, and server-agent settings.
---

# HetrixTools Provider

The HetrixTools provider manages resources exposed by the HetrixTools API.

## Example Usage

```terraform
terraform {
  required_providers {
    hetrixtools = {
      source  = "xaf/hetrixtools"
      version = "~> 0.1"
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

## Schema

### Optional

- `api_token` (String, Sensitive) HetrixTools API bearer token. Can also be set with `HETRIXTOOLS_API_TOKEN`.
- `base_url` (String) HetrixTools API root URL. The provider appends `/v2` and `/v3`. Can also be set with `HETRIXTOOLS_BASE_URL`.
- `base_url_v2` (String) HetrixTools v2 API base URL. Overrides `base_url` for token-path endpoints. Can also be set with `HETRIXTOOLS_BASE_URL_V2`.
- `base_url_v3` (String) HetrixTools v3 API base URL. Overrides `base_url` for REST endpoints. Can also be set with `HETRIXTOOLS_BASE_URL_V3`.
