package hetrixtools

import "context"

// GetAccountLimits returns the current account-level HetrixTools limits.
func (c *Client) GetAccountLimits(ctx context.Context) (any, error) {
	body, err := c.getEndpoint(ctx, "/account/limits", nil)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
