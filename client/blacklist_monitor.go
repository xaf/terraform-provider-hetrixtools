package hetrixtools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CreateBlacklistMonitor creates a HetrixTools blacklist monitor using the
// documented v2 blacklist add endpoint:
//
//   - https://docs.hetrixtools.com/api-add-blacklist-monitor/
func (c *Client) CreateBlacklistMonitor(ctx context.Context, request BlacklistMonitorRequest) (*ActionResponse, error) {
	body, err := c.doV2Form(ctx, "/blacklist/add/", request.form())
	if err != nil {
		return nil, err
	}
	c.clearMonitorCaches()
	return decodeActionResponse(body)
}

// UpdateBlacklistMonitor updates a HetrixTools blacklist monitor using the
// documented v2 blacklist edit endpoint:
//
//   - https://docs.hetrixtools.com/api-edit-blacklist-monitor/
func (c *Client) UpdateBlacklistMonitor(ctx context.Context, request BlacklistMonitorRequest) (*ActionResponse, error) {
	body, err := c.doV2Form(ctx, "/blacklist/edit/", request.form())
	if err != nil {
		return nil, err
	}
	c.clearMonitorCaches()
	return decodeActionResponse(body)
}

// UpsertBlacklistMonitor updates an existing blacklist monitor by target or creates it when absent.
func (c *Client) UpsertBlacklistMonitor(ctx context.Context, request BlacklistMonitorRequest) (*ActionResponse, error) {
	existing, err := c.GetBlacklistMonitor(ctx, request.Target)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return c.CreateBlacklistMonitor(ctx, request)
	}
	return c.UpdateBlacklistMonitor(ctx, request)
}

// DeleteBlacklistMonitor deletes a HetrixTools blacklist monitor by target using
// the documented v2 blacklist delete endpoint:
//
//   - https://docs.hetrixtools.com/api-delete-blacklist-monitor/
func (c *Client) DeleteBlacklistMonitor(ctx context.Context, target string) error {
	_, err := c.doV2Form(ctx, "/blacklist/delete/", url.Values{"target": {target}})
	if err == nil {
		c.clearMonitorCaches()
	}
	return err
}

// ListBlacklistMonitors returns blacklist monitors matching query filters.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1blacklist-monitors/get
func (c *Client) ListBlacklistMonitors(ctx context.Context, query map[string]string) (*BlacklistMonitorsResponse, error) {
	var response BlacklistMonitorsResponse
	if err := c.getJSON(ctx, "/blacklist-monitors", query, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetBlacklistMonitor finds a blacklist monitor by exact target.
func (c *Client) GetBlacklistMonitor(ctx context.Context, target string) (*BlacklistMonitor, error) {
	monitors, err := c.cachedBlacklistMonitors(ctx)
	if err != nil {
		return nil, err
	}
	for _, monitor := range monitors {
		if monitor.Target == target {
			monitor := monitor
			return &monitor, nil
		}
	}
	return nil, nil
}

func (c *Client) cachedBlacklistMonitors(ctx context.Context) ([]BlacklistMonitor, error) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if c.blMonitors != nil {
		return c.blMonitors, nil
	}

	var monitors []BlacklistMonitor
	for page := 1; ; page++ {
		response, err := c.ListBlacklistMonitors(ctx, map[string]string{"page": fmt.Sprint(page), "per_page": "100"})
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, response.BlacklistMonitors...)
		if response.Meta.Pagination.Next == nil || page >= response.Meta.Pagination.Last {
			c.blMonitors = monitors
			return c.blMonitors, nil
		}
	}
}

// GetBlacklistMonitorReport returns the report for a blacklist monitor
// identifier as a decoded JSON value, typically a map[string]any. Query keys are
// passed through to the HetrixTools v3 report endpoint. Source-of-truth API
// docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1blacklist-monitors~1{identifier}~1report/get
func (c *Client) GetBlacklistMonitorReport(ctx context.Context, identifier string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/blacklist-monitors/"+identifier+"/report", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}

// CheckBlacklistIPv4 runs a one-off IPv4 blacklist check using the documented
// v2 blacklist check endpoint. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/blacklist-check-api/
func (c *Client) CheckBlacklistIPv4(ctx context.Context, ipAddress string) (*BlacklistCheckResult, error) {
	return c.checkBlacklist(ctx, "ipv4", ipAddress)
}

// CheckBlacklistDomain runs a one-off domain blacklist check using the
// documented v2 blacklist check endpoint. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/blacklist-check-api/
func (c *Client) CheckBlacklistDomain(ctx context.Context, domain string) (*BlacklistCheckResult, error) {
	return c.checkBlacklist(ctx, "domain", domain)
}

func (c *Client) checkBlacklist(ctx context.Context, kind string, target string) (*BlacklistCheckResult, error) {
	body, err := c.doV2JSON(ctx, http.MethodGet, "/blacklist-check/"+kind+"/"+url.PathEscape(target)+"/", nil)
	if err != nil {
		return nil, err
	}
	var result BlacklistCheckResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	result.RawJSON = append(result.RawJSON[:0], body...)
	return &result, nil
}
