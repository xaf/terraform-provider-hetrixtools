package hetrixtools

import "context"

// GetAccountLimits returns the current account-level HetrixTools limits as a
// decoded JSON value, typically a map[string]any. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1account~1limits/get
func (c *Client) GetAccountLimits(ctx context.Context) (any, error) {
	body, err := c.getEndpoint(ctx, "/account/limits", nil)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
