package hetrixtools

import "context"

// GetAccountLimits returns the current account-level HetrixTools limits as a
// decoded JSON value, typically a map[string]any. The endpoint is documented in
// the HetrixTools v3 API reference.
func (c *Client) GetAccountLimits(ctx context.Context) (any, error) {
	body, err := c.getEndpoint(ctx, "/account/limits", nil)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
