package hetrixtools

import "context"

// ListContactLists returns HetrixTools contact lists matching query filters as a
// decoded JSON value, typically a map[string]any. Supported query keys are the
// HetrixTools v3 contact-list endpoint's query parameters, such as page and
// per_page. Source-of-truth API docs:
//
//   - https://docs.hetrixtools.com/api/v3/#/paths/~1contact-lists/get
func (c *Client) ListContactLists(ctx context.Context, query map[string]string) (any, error) {
	body, err := c.getEndpoint(ctx, "/contact-lists", query)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}
