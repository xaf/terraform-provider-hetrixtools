package hetrixtools

import (
	"context"
	"fmt"
	"net/http"
)

// CreateUptimeMonitor creates a HetrixTools uptime monitor using the documented
// v2 uptime add endpoint:
//
//   - https://docs.hetrixtools.com/api-add-website-ping-service-smtp-uptime-monitor/
//   - https://docs.hetrixtools.com/api-add-server-agent-uptime-monitor-heartbeat-uptime-monitor/
func (c *Client) CreateUptimeMonitor(ctx context.Context, request UptimeMonitorRequest) (*ActionResponse, error) {
	request.MID = ""
	body, err := c.doV2JSON(ctx, http.MethodPost, "/uptime/add/", request)
	if err != nil {
		return nil, err
	}
	c.clearMonitorCaches()
	return decodeActionResponse(body)
}

// UpdateUptimeMonitor updates a HetrixTools uptime monitor using the documented
// v2 uptime add endpoint with MID set, as described by HetrixTools:
//
//   - https://docs.hetrixtools.com/api-add-website-ping-service-smtp-uptime-monitor/
//   - https://docs.hetrixtools.com/api-add-server-agent-uptime-monitor-heartbeat-uptime-monitor/
func (c *Client) UpdateUptimeMonitor(ctx context.Context, request UptimeMonitorRequest) (*ActionResponse, error) {
	body, err := c.doV2JSON(ctx, http.MethodPost, "/uptime/add/", request)
	if err != nil {
		return nil, err
	}
	c.clearMonitorCaches()
	return decodeActionResponse(body)
}

// UpsertUptimeMonitor updates an uptime monitor when MID is set, otherwise it creates one.
func (c *Client) UpsertUptimeMonitor(ctx context.Context, request UptimeMonitorRequest) (*ActionResponse, error) {
	if request.MID == "" {
		return c.CreateUptimeMonitor(ctx, request)
	}
	return c.UpdateUptimeMonitor(ctx, request)
}

// DeleteUptimeMonitor deletes a HetrixTools uptime monitor by monitor ID using
// the documented v2 uptime delete endpoint:
//
//   - https://docs.hetrixtools.com/api-delete-uptime-monitor/
func (c *Client) DeleteUptimeMonitor(ctx context.Context, monitorID string) error {
	_, err := c.doV2JSON(ctx, http.MethodPost, "/uptime/delete/", map[string]string{"MID": monitorID})
	if err == nil {
		c.clearMonitorCaches()
	}
	return err
}

// ListUptimeMonitors returns uptime monitors matching query filters.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1uptime-monitors/get
func (c *Client) ListUptimeMonitors(ctx context.Context, query map[string]string) (*UptimeMonitorsResponse, error) {
	var response UptimeMonitorsResponse
	if err := c.getJSON(ctx, "/uptime-monitors", query, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUptimeMonitor finds an uptime monitor by monitor ID.
func (c *Client) GetUptimeMonitor(ctx context.Context, monitorID string) (*UptimeMonitor, error) {
	if monitorID == "" {
		return nil, nil
	}
	monitors, err := c.cachedUptimeMonitors(ctx)
	if err != nil {
		return nil, err
	}
	for _, monitor := range monitors {
		if monitor.ID == monitorID {
			monitor := monitor
			return &monitor, nil
		}
	}
	return nil, nil
}

func (c *Client) cachedUptimeMonitors(ctx context.Context) ([]UptimeMonitor, error) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	if c.uptimeMonitors != nil {
		return c.uptimeMonitors, nil
	}

	var monitors []UptimeMonitor
	for page := 1; ; page++ {
		response, err := c.ListUptimeMonitors(ctx, map[string]string{"page": fmt.Sprint(page), "per_page": "100"})
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, response.UptimeMonitors...)
		if response.Meta.Pagination.Next == nil || page >= response.Meta.Pagination.Last {
			c.uptimeMonitors = monitors
			return c.uptimeMonitors, nil
		}
	}
}

// GetUptimeMonitorReport returns a report for an uptime monitor as a decoded
// JSON value, typically a map[string]any. Query keys are passed through to the
// HetrixTools v3 report endpoint. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1uptime-monitors~1{monitor_id}~1report/get
func (c *Client) GetUptimeMonitorReport(ctx context.Context, monitorID string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/report", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}

// ListUptimeMonitorDowntimes returns downtime entries for an uptime monitor as a
// decoded JSON value, typically a map[string]any. Query keys are passed through
// to the HetrixTools v3 downtime endpoint. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1uptime-monitors~1{monitor_id}~1downtimes/get
func (c *Client) ListUptimeMonitorDowntimes(ctx context.Context, monitorID string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/downtimes", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}

// GetUptimeMonitorLocationFailLog returns location failure logs for an uptime
// monitor as a decoded JSON value, typically a map[string]any. Query keys are
// passed through to the HetrixTools v3 location-fail-log endpoint.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1uptime-monitors~1{monitor_id}~1location-fail-log/get
func (c *Client) GetUptimeMonitorLocationFailLog(ctx context.Context, monitorID string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/location-fail-log", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
