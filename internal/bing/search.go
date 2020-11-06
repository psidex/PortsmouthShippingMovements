package bing

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const apiUrl = "https://api.bing.microsoft.com/v7.0/images/search?count=1&imageType=Photo&license=Share&q="

var RequestRateError = errors.New("request rate limit exceeded")

// imageSearch holds the data needed from a Bing image search API request.
type imageSearch struct {
	Value []struct {
		ThumbnailURL string `json:"thumbnailUrl"`
	} `json:"value"`
	Error struct {
		Code int `json:"code"`
	} `json:"error"`
}

// ImageSearchApi contains methods for interacting with the Bing image search API.
type ImageSearchApi struct {
	client *http.Client
	apiKey string
}

// NewImageSearchApi creates a new ImageSearchApi.
func NewImageSearchApi(client *http.Client, apiKey string) ImageSearchApi {
	return ImageSearchApi{client: client, apiKey: apiKey}
}

// SearchForImage attempts to find a thumbnail image URL for the given query.
func (i ImageSearchApi) SearchForImage(query string) (string, error) {
	queryUrl := apiUrl + url.QueryEscape(query)
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

	if data.Error.Code == 429 {
		return "", RequestRateError
	}
	if len(data.Value) <= 0 {
		// IThis will also trigger if we hit a rate limit.
		return "", errors.New("no images found")
	}
	return data.Value[0].ThumbnailURL, nil
}
