package bing

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// imageSearch holds the data needed from a Bing image search API request.
type imageSearch struct {
	Value []struct {
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"value"`
}

// ImageSearchApi contains methods for interacting with the Bing image search API.
type ImageSearchApi struct {
	client  *http.Client
	apiKey  string
	baseUrl string // The base for creating an Image API request. Must end with "&q=".
}

// NewImageSearchApi creates a new ImageSearchApi.
func NewImageSearchApi(client *http.Client, apiKey string, baseUrl string) ImageSearchApi {
	return ImageSearchApi{client: client, apiKey: apiKey, baseUrl: baseUrl}
}

// SearchForImage attempts to find a thumbnail image URL for the given query.
func (i ImageSearchApi) SearchForImage(query string) (string, error) {
	queryUrl := i.baseUrl + url.QueryEscape(query)
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", i.apiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data imageSearch
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return "", err
	}

	if len(data.Value) <= 0 {
		// IThis will also trigger if we hit a rate limit.
		return "", errors.New("no images found")
	}
	return data.Value[0].ThumbnailURL, nil
}
