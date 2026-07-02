package hetrixtools

import (
	"context"
	"fmt"
	"net/http"
)

// CreateUptimeMonitor creates a HetrixTools uptime monitor.
func (c *Client) CreateUptimeMonitor(ctx context.Context, request UptimeMonitorRequest) (*ActionResponse, error) {
	request.MID = ""
	body, err := c.doV2JSON(ctx, http.MethodPost, "/uptime/add/", request)
	if err != nil {
		return nil, err
	}
	return decodeActionResponse(body)
}

// UpdateUptimeMonitor updates a HetrixTools uptime monitor.
func (c *Client) UpdateUptimeMonitor(ctx context.Context, request UptimeMonitorRequest) (*ActionResponse, error) {
	body, err := c.doV2JSON(ctx, http.MethodPost, "/uptime/add/", request)
	if err != nil {
		return nil, err
	}
	return decodeActionResponse(body)
}

// UpsertUptimeMonitor updates an uptime monitor when MID is set, otherwise it creates one.
func (c *Client) UpsertUptimeMonitor(ctx context.Context, request UptimeMonitorRequest) (*ActionResponse, error) {
	if request.MID == "" {
		return c.CreateUptimeMonitor(ctx, request)
	}
	return c.UpdateUptimeMonitor(ctx, request)
}

// DeleteUptimeMonitor deletes a HetrixTools uptime monitor by monitor ID.
func (c *Client) DeleteUptimeMonitor(ctx context.Context, monitorID string) error {
	_, err := c.doV2JSON(ctx, http.MethodPost, "/uptime/delete/", map[string]string{"MID": monitorID})
	return err
}

// ListUptimeMonitors returns uptime monitors matching query filters.
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
	for page := 1; ; page++ {
		response, err := c.ListUptimeMonitors(ctx, map[string]string{"page": fmt.Sprint(page), "per_page": "100", "id": monitorID})
		if err != nil {
			return nil, err
		}
		for _, monitor := range response.UptimeMonitors {
			if monitor.ID == monitorID {
				return &monitor, nil
			}
		}
		if response.Meta.Pagination.Next == nil || page >= response.Meta.Pagination.Last {
			return nil, nil
		}
	}
}

// GetUptimeMonitorReport returns a report for an uptime monitor.
func (c *Client) GetUptimeMonitorReport(ctx context.Context, monitorID string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/report", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}

// ListUptimeMonitorDowntimes returns downtime entries for an uptime monitor.
func (c *Client) ListUptimeMonitorDowntimes(ctx context.Context, monitorID string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/downtimes", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}

// GetUptimeMonitorLocationFailLog returns location failure logs for an uptime monitor.
func (c *Client) GetUptimeMonitorLocationFailLog(ctx context.Context, monitorID string, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/location-fail-log", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
