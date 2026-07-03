package hetrixtools

import (
	"context"
	"fmt"
)

// CreateScheduledMaintenance creates a scheduled maintenance window.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1schedule-maintenance/post
func (c *Client) CreateScheduledMaintenance(ctx context.Context, request ScheduledMaintenanceRequest) (*ScheduledMaintenance, error) {
	var created ScheduledMaintenance
	if err := c.postJSON(ctx, "/schedule-maintenance", request, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// DeleteScheduledMaintenance deletes a scheduled maintenance window by ID.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1schedule-maintenance~1{schedule_maintenance_id}/delete
func (c *Client) DeleteScheduledMaintenance(ctx context.Context, id string) error {
	return c.deleteJSON(ctx, "/schedule-maintenance/"+id, nil)
}

// ListScheduledMaintenances returns scheduled maintenance windows matching
// query filters. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1schedule-maintenance/get
func (c *Client) ListScheduledMaintenances(ctx context.Context, query map[string]string) (*ScheduledMaintenancesResponse, error) {
	var response ScheduledMaintenancesResponse
	if err := c.getJSON(ctx, "/schedule-maintenance", query, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetScheduledMaintenance finds a scheduled maintenance window by ID.
func (c *Client) GetScheduledMaintenance(ctx context.Context, id string, monitorID string) (*ScheduledMaintenance, error) {
	for page := 1; ; page++ {
		query := map[string]string{"page": fmt.Sprint(page), "per_page": "200"}
		if monitorID != "" {
			query["monitor_id"] = monitorID
		}
		response, err := c.ListScheduledMaintenances(ctx, query)
		if err != nil {
			return nil, err
		}
		for _, maintenance := range response.ScheduledMaintenances {
			if maintenance.ID == id {
				return &maintenance, nil
			}
		}
		if response.Meta.Pagination.Next == nil || page >= response.Meta.Pagination.Last {
			return nil, nil
		}
	}
}
