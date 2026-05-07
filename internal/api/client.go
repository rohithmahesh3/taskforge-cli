package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rohithmahesh3/taskforge-cli/internal/config"
)

const (
	DefaultTimeout = 30 * time.Second
	APIVersion     = "v1"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
	Workspace  string
}

type Pagination struct {
	Count  int     `json:"count"`
	Next   *string `json:"next"`
	Prev   *string `json:"prev"`
	Offset *int    `json:"offset,omitempty"`
}

type Response struct {
	Pagination
	Results json.RawMessage `json:"results"`
}

type APIError struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	RawBody    string          `json:"raw_body"`
	Conflicts  json.RawMessage `json:"conflicts,omitempty"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.RawBody)
}

func NewClient() (*Client, error) {
	apiKey, err := config.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Run 'taskforge auth login' first")
	}

	workspace := config.Cfg.DefaultWorkspace
	if workspace == "" {
		return nil, fmt.Errorf("no default workspace set. Use --workspace flag or set default workspace")
	}

	baseURL := config.Cfg.APIHost
	if baseURL == "" {
		return nil, fmt.Errorf("no API host configured. Run 'taskforge auth login' to configure")
	}

	return &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Workspace:  workspace,
	}, nil
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s/api/%s%s", c.BaseURL, APIVersion, path)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	body, err := c.DoRaw(req)
	if err != nil {
		return err
	}

	if v != nil && len(body) > 0 {
		return json.Unmarshal(body, v)
	}

	return nil
}

func (c *Client) Get(path string, query url.Values, v interface{}) error {
	if query != nil {
		path = path + "?" + query.Encode()
	}

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) DoRaw(req *http.Request) ([]byte, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    string(body),
		}
		var parsed struct {
			Error     string          `json:"error"`
			Message   string          `json:"message"`
			Conflicts json.RawMessage `json:"conflicts"`
		}
		if err := json.Unmarshal(body, &parsed); err == nil {
			if parsed.Error != "" {
				apiErr.Message = parsed.Error
			} else if parsed.Message != "" {
				apiErr.Message = parsed.Message
			}
			apiErr.Conflicts = parsed.Conflicts
		}
		return nil, apiErr
	}

	return body, nil
}

func (c *Client) GetRaw(path string, query url.Values) ([]byte, error) {
	if query != nil {
		path = path + "?" + query.Encode()
	}

	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.DoRaw(req)
}

func (c *Client) Post(path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest("POST", path, body)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) Patch(path string, body interface{}, v interface{}) error {
	req, err := c.NewRequest("PATCH", path, body)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) Delete(path string) error {
	req, err := c.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return c.Do(req, nil)
}

func (c *Client) SetWorkspace(workspace string) {
	c.Workspace = workspace
}

func unmarshalListResponse[T any](body []byte) ([]T, error) {
	var wrapped struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && wrapped.Results != nil {
		var items []T
		if err := json.Unmarshal(wrapped.Results, &items); err != nil {
			return nil, err
		}
		return items, nil
	}

	var items []T
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func isAPIStatusError(err error, statusCode int) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == statusCode
	}
	return strings.Contains(err.Error(), fmt.Sprintf("status %d", statusCode))
}
