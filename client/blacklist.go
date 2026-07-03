package hetrixtools

import "context"

// ListBlacklists returns the blacklist providers monitored by HetrixTools as a
// decoded JSON value, typically a map[string]any. The endpoint is documented in
// the HetrixTools v3 API reference.
func (c *Client) ListBlacklists(ctx context.Context) (any, error) {
	body, err := c.getEndpoint(ctx, "/blacklists", nil)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
