package apiClient

import (
	"SongLibrary/pkg/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	logrus.Debug("Creating new API client with base URL: " + baseURL)
	return &Client{
		BaseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) FetchSongDetails(group, songName string) (*model.SongAPIResponse, error) {
	logrus.Debugf("Fetching song details for group: %s, song: %s", group, songName)

	if group == "" || songName == "" {
		logrus.Warn("group and songName must not be empty")
		return nil, fmt.Errorf("group and songName must not be empty")
	}

	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.BaseURL, group, songName)
	logrus.Debugf("Requesting URL: %s", url)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		logrus.Errorf("Failed to send request: %s", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Unexpected response status: %s", resp.Status)
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var apiResponse model.SongAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		logrus.Errorf("Failed to decode API response: %s", err)
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	logrus.Debug("Successfully fetched song details")
	return &apiResponse, nil
}
