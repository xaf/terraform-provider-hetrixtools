package hetrixtools

import "context"

// ListBlacklists returns the blacklist providers monitored by HetrixTools.
func (c *Client) ListBlacklists(ctx context.Context) (any, error) {
	body, err := c.getEndpoint(ctx, "/blacklists", nil)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
