package hetrixtools

import "context"

// AttachServerAgent attaches a server agent to an uptime monitor.
func (c *Client) AttachServerAgent(ctx context.Context, monitorID string) (*ServerAgentResponse, error) {
	var result ServerAgentResponse
	if err := c.postJSON(ctx, "/uptime-monitors/"+monitorID+"/server-agent", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetServerAgent returns the server agent attached to an uptime monitor.
func (c *Client) GetServerAgent(ctx context.Context, monitorID string) (*ServerAgentResponse, error) {
	var result ServerAgentResponse
	if err := c.getJSON(ctx, "/uptime-monitors/"+monitorID+"/server-agent", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DetachServerAgent detaches a server agent from an uptime monitor.
func (c *Client) DetachServerAgent(ctx context.Context, monitorID string) error {
	return c.deleteJSON(ctx, "/uptime-monitors/"+monitorID+"/server-agent", nil)
}

// GetServerAgentWarningPolicies returns server-agent warning policies for an
// uptime monitor as a decoded JSON value, typically a map[string]any. The shape
// mirrors the HetrixTools v3 warning-policies payload.
func (c *Client) GetServerAgentWarningPolicies(ctx context.Context, monitorID string) (any, error) {
	body, err := c.getEndpoint(ctx, "/uptime-monitors/"+monitorID+"/server-agent/warning-policies", nil)
	if err != nil {
		return nil, err
	}
	return decodeUntypedJSON(body)
}

// UpdateServerAgentWarningPolicies replaces server-agent warning policies for an uptime monitor.
func (c *Client) UpdateServerAgentWarningPolicies(ctx context.Context, monitorID string, payload any) error {
	return c.putJSON(ctx, "/uptime-monitors/"+monitorID+"/server-agent/warning-policies", payload, nil)
}
