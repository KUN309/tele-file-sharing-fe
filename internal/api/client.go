package api

import (
	"io"
	"fmt"
	"path"
	"bytes"
	"net/http"
	"encoding/json"
)

type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTP:    &http.Client{},
	}
}

// From Go 1.18+, 'any' is interface{} alias
func (c *Client) request(method, p string, body any) ([]byte, int, error) {
	var buf io.Reader

	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, c.BaseURL+p, buf)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	return data, resp.StatusCode, nil
}

// ---------------- USER ----------------

func (c *Client) GetMe() (*User, error) {
	b, status, err := c.request("GET", "/api/me", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("status %d: %s", status, string(b))
	}

	var res MeResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// ---------------- FILES ----------------

func (c *Client) ListFiles() ([]FileItem, error) {
	b, status, err := c.request("GET", "/api/v1/files", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("status %d: %s", status, string(b))
	}

	var files []FileItem
	json.Unmarshal(b, &files)
	return files, nil
}

func (c *Client) UploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	b, status, err := c.request("POST", "/api/v1/files", req)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("tải lên thất bại: %s", string(b))
	}

	var res UploadFileResponse
	json.Unmarshal(b, &res)
	return &res, nil
}

// ---------------- SHARES ----------------

func (c *Client) CreateShare(req ShareCreateRequest) (*ShareItem, error) {
	b, status, err := c.request("POST", "/api/v1/shares", req)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("chia sẻ thất bại: %s", string(b))
	}

	var share ShareItem
	json.Unmarshal(b, &share)
	return &share, nil
}

func (c *Client) ListShares() ([]ShareItem, error) {
	b, status, err := c.request("GET", "/api/v1/shares", nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("status %d: %s", status, string(b))
	}

	var shares []ShareItem
	json.Unmarshal(b, &shares)
	return shares, nil
}

func (c *Client) RevokeShare(id int) (*RevokeResponse, error) {
	url := path.Join("/api/v1/shares", fmt.Sprint(id), "revoke")
	b, status, err := c.request("POST", url, nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("thu hồi thất bại: %s", string(b))
	}

	var res RevokeResponse
	json.Unmarshal(b, &res)
	return &res, nil
}

func (c *Client) AuthorizeShare(id int, password string) error {
	url := path.Join("/api/v1/shares", fmt.Sprint(id), "authorize")
	req := PasswordAuthorizeRequest{Password: password}

	b, status, err := c.request("POST", url, req)
	if err != nil {
		return err
	}
	if status >= 400 {
		return fmt.Errorf("cấp quyền thất bại: %s", string(b))
	}
	return nil
}

func (c *Client) GetShareMetadata(id int) (*ShareMetadata, error) {
	url := path.Join("/api/v1/shares", fmt.Sprint(id))
	b, status, err := c.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("metadata error: %s", string(b))
	}

	var m ShareMetadata
	json.Unmarshal(b, &m)
	return &m, nil
}

func (c *Client) DownloadShare(id int) ([]byte, error) {
	url := path.Join("/api/v1/shares", fmt.Sprint(id), "download")

	// lấy file binary
	req, _ := http.NewRequest("GET", c.BaseURL+url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}