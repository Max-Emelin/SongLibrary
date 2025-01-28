package apiClient

import (
	"SongLibrary/pkg/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) FetchSongDetails(group, songName string) (*model.SongAPIResponse, error) {
	if group == "" || songName == "" {
		return nil, fmt.Errorf("group and songName must not be empty")
	}

	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.BaseURL, group, songName)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var apiResponse model.SongAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return &apiResponse, nil
}
