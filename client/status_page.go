package hetrixtools

import (
	"context"
	"fmt"
)

// ListStatusPages returns status pages matching query filters.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1status-pages/get
func (c *Client) ListStatusPages(ctx context.Context, query map[string]string) (*StatusPagesResponse, error) {
	var response StatusPagesResponse
	if err := c.getJSON(ctx, "/status-pages", query, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetStatusPage finds a status page by ID using ListStatusPages; see
// ListStatusPages for the source API reference.
func (c *Client) GetStatusPage(ctx context.Context, id string) (*StatusPage, error) {
	for page := 1; ; page++ {
		response, err := c.ListStatusPages(ctx, map[string]string{"page": fmt.Sprint(page), "per_page": "100"})
		if err != nil {
			return nil, err
		}
		for _, statusPage := range response.StatusPages {
			if statusPage.ID == id {
				return &statusPage, nil
			}
		}
		if response.Meta.Pagination.Next == nil || page >= response.Meta.Pagination.Last {
			return nil, nil
		}
	}
}

// AddStatusPageMonitors adds uptime monitors to a status page.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1status-pages~1{status_page_id}~1monitors/post
func (c *Client) AddStatusPageMonitors(ctx context.Context, statusPageID string, monitorIDs []string) error {
	return c.postJSON(ctx, "/status-pages/"+statusPageID+"/monitors", monitorIDs, nil)
}

// RemoveStatusPageMonitors removes uptime monitors from a status page.
// Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1status-pages~1{status_page_id}~1monitors/delete
func (c *Client) RemoveStatusPageMonitors(ctx context.Context, statusPageID string, monitorIDs []string) error {
	return c.deleteJSON(ctx, "/status-pages/"+statusPageID+"/monitors", monitorIDs)
}
