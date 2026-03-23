package alfabank

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const defaultBaseURL = "https://alfa.rbsuat.com/payment/rest"

// Config configures Client construction.
type Config struct {
	BaseURL    string
	UserName   string
	Password   string
	HTTPClient *http.Client
}

// Client is a minimal Alfa acquiring REST API client.
type Client struct {
	baseURL    string
	userName   string
	password   string
	httpClient *http.Client
}

// NewClient constructs a client for the Alfa REST API.
func NewClient(cfg Config) *Client {
	baseURL := cfg.BaseURL
	if strings.TrimSpace(baseURL) == "" {
		baseURL = defaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		userName:   cfg.UserName,
		password:   cfg.Password,
		httpClient: httpClient,
	}
}

func (c *Client) postForm(ctx context.Context, methodPath string, req any, resp any) error {
	values, err := encodeForm(req)
	if err != nil {
		return err
	}

	values.Set("userName", c.userName)
	values.Set("password", c.password)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+strings.TrimLeft(methodPath, "/"), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", httpResp.Status)
	}

	if resp == nil {
		return nil
	}

	decoder := json.NewDecoder(httpResp.Body)
	if err := decoder.Decode(resp); err != nil {
		return err
	}

	return nil
}