package hetrixtools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default root URL for HetrixTools API endpoints.
	DefaultBaseURL = "https://api.hetrixtools.com"
	// DefaultV3BaseURL is the default base URL for HetrixTools REST endpoints.
	DefaultV3BaseURL = DefaultBaseURL + "/v3"
	// DefaultV2BaseURL is the default base URL for older HetrixTools token-path endpoints.
	DefaultV2BaseURL = DefaultBaseURL + "/v2"
)

// Client calls HetrixTools APIs through semantic resource methods.
type Client struct {
	v3BaseURL string
	v2BaseURL string
	token     string
	http      *http.Client
}

type Option func(*Client)

// WithHTTPClient configures the HTTP client used for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.http = httpClient
		}
	}
}

// WithV2BaseURL overrides the base URL used for older token-path endpoints.
func WithV2BaseURL(baseURL string) Option {
	return func(c *Client) {
		if strings.TrimSpace(baseURL) != "" {
			c.v2BaseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// WithV3BaseURL overrides the base URL used for REST endpoints.
func WithV3BaseURL(baseURL string) Option {
	return func(c *Client) {
		if strings.TrimSpace(baseURL) != "" {
			c.v3BaseURL = strings.TrimRight(baseURL, "/")
		}
	}
}

// NewClient returns a client configured with the default HetrixTools base URLs.
func NewClient(token string, options ...Option) *Client {
	return NewClientWithBaseURL(DefaultBaseURL, token, options...)
}

// NewClientWithBaseURL returns a client configured with a custom API root URL.
// For compatibility, baseURL may also be a versioned /v2 or /v3 URL.
func NewClientWithBaseURL(baseURL string, token string, options ...Option) *Client {
	v2BaseURL, v3BaseURL := versionedBaseURLs(baseURL)

	c := &Client{
		v3BaseURL: v3BaseURL,
		v2BaseURL: v2BaseURL,
		token:     token,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Client) getJSON(ctx context.Context, path string, query map[string]string, out any) error {
	body, err := c.doV3(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(body, out)
}

func (c *Client) postJSON(ctx context.Context, path string, body any, out any) error {
	responseBody, err := c.doV3(ctx, http.MethodPost, path, nil, body)
	if err != nil {
		return err
	}
	if out == nil || len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

func (c *Client) putJSON(ctx context.Context, path string, body any, out any) error {
	responseBody, err := c.doV3(ctx, http.MethodPut, path, nil, body)
	if err != nil {
		return err
	}
	if out == nil || len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

func (c *Client) deleteJSON(ctx context.Context, path string, body any) error {
	_, err := c.doV3(ctx, http.MethodDelete, path, nil, body)
	return err
}

func (c *Client) getEndpoint(ctx context.Context, path string, query map[string]string) ([]byte, error) {
	return c.doV3(ctx, http.MethodGet, path, query, nil)
}

func (c *Client) doV3(ctx context.Context, method string, path string, query map[string]string, body any) ([]byte, error) {
	return c.do(ctx, c.v3BaseURL, method, path, query, body, true)
}

func (c *Client) doV2JSON(ctx context.Context, method string, path string, body any) ([]byte, error) {
	return c.do(ctx, c.v2BaseURL, method, c.v2Path(path), nil, body, false)
}

func (c *Client) doV2Form(ctx context.Context, path string, values url.Values) ([]byte, error) {
	return c.doForm(ctx, c.v2BaseURL, c.v2Path(path), values)
}

func decodeActionResponse(body []byte) (*ActionResponse, error) {
	var result ActionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Status == "ERROR" || result.ErrorMessage != "" {
		return &result, Error{Response: &result}
	}
	return &result, nil
}

func (c *Client) do(ctx context.Context, baseURL string, method string, path string, query map[string]string, body any, bearerAuth bool) ([]byte, error) {
	requestURL, err := requestURL(baseURL, path, query)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerAuth && c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.doRequest(req)
}

func (c *Client) doForm(ctx context.Context, baseURL string, path string, values url.Values) ([]byte, error) {
	requestURL, err := requestURL(baseURL, path, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.doRequest(req)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, Error{StatusCode: resp.StatusCode, Body: string(responseBody)}
	}

	return responseBody, nil
}

func (c *Client) v2Path(path string) string {
	return "/" + url.PathEscape(c.token) + "/" + strings.TrimLeft(path, "/")
}

func requestURL(baseURL string, path string, query map[string]string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return "", fmt.Errorf("endpoint path must be relative, got %q", path)
	}

	u, err := url.Parse(strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/"))
	if err != nil {
		return "", err
	}

	values := u.Query()
	for key, value := range query {
		if value != "" {
			values.Set(key, value)
		}
	}
	u.RawQuery = values.Encode()

	return u.String(), nil
}

func versionedBaseURLs(baseURL string) (string, string) {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	if strings.HasSuffix(baseURL, "/v2") || strings.HasSuffix(baseURL, "/v3") {
		baseURL = baseURL[:len(baseURL)-3]
	}
	return baseURL + "/v2", baseURL + "/v3"
}

// Error is the common error type returned by the client for HetrixTools API failures.
type Error struct {
	StatusCode int
	Body       string
	Response   *ActionResponse
}

// Error returns a human-readable HetrixTools API error string.
func (e Error) Error() string {
	if e.Response != nil {
		if e.Response.ErrorMessage != "" {
			return e.Response.ErrorMessage
		}
		return "hetrixtools API action failed"
	}
	if e.Body == "" {
		return fmt.Sprintf("hetrixtools API returned HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("hetrixtools API returned HTTP %d: %s", e.StatusCode, e.Body)
}

// IsNotFound reports whether err is a HetrixTools 404 response.
func IsNotFound(err error) bool {
	var apiErr Error
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}
