// Package provider implements the Terraform Plugin Framework provider for
// HetrixTools resources and data sources.
//
// API behavior such as version selection, legacy v2 payload translation, monitor
// cache invalidation, and request pacing is delegated to the reusable client
// package.
package provider
