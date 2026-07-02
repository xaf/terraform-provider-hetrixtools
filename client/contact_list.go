package hetrixtools

import "context"

// ListContactLists returns HetrixTools contact lists matching query filters.
func (c *Client) ListContactLists(ctx context.Context, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/contact-lists", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
