terraform {
  required_providers {
    hetrixtools = {
      source = "github.com/xaf/go-hetrixtools"
    }
  }
}

provider "hetrixtools" {}

data "hetrixtools_account_limits" "current" {}

data "hetrixtools_uptime_monitors" "all" {
  per_page = "200"
}
