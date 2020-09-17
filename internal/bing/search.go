package bing

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// baseUrl is the base route for the image search API.
const baseUrl = "https://api.cognitive.microsoft.com/bing/v7.0/images/search?q="

// imageSearch represents the API response from a BingApi image search API request.
// https://mholt.github.io/json-to-go/
type imageSearch struct {
	Type            string `json:"_type"`
	Instrumentation struct {
		Type string `json:"_type"`
	} `json:"instrumentation"`
	ReadLink     string `json:"readLink"`
	WebSearchURL string `json:"webSearchUrl"`
	QueryContext struct {
		OriginalQuery           string `json:"originalQuery"`
		AlterationDisplayQuery  string `json:"alterationDisplayQuery"`
		AlterationOverrideQuery string `json:"alterationOverrideQuery"`
		AlterationMethod        string `json:"alterationMethod"`
		AlterationType          string `json:"alterationType"`
	} `json:"queryContext"`
	TotalEstimatedMatches int `json:"totalEstimatedMatches"`
	NextOffset            int `json:"nextOffset"`
	CurrentOffset         int `json:"currentOffset"`
	Value                 []struct {
		WebSearchURL       string    `json:"webSearchUrl"`
		Name               string    `json:"name"`
		ThumbnailURL       string    `json:"thumbnailUrl"`
		DatePublished      time.Time `json:"datePublished"`
		IsFamilyFriendly   bool      `json:"isFamilyFriendly"`
		ContentURL         string    `json:"contentUrl"`
		HostPageURL        string    `json:"hostPageUrl"`
		ContentSize        string    `json:"contentSize"`
		EncodingFormat     string    `json:"encodingFormat"`
		HostPageDisplayURL string    `json:"hostPageDisplayUrl"`
		Width              int       `json:"width"`
		Height             int       `json:"height"`
		Thumbnail          struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"thumbnail"`
		ImageInsightsToken string `json:"imageInsightsToken"`
		InsightsMetadata   struct {
			RecipeSourcesCount  int `json:"recipeSourcesCount"`
			PagesIncludingCount int `json:"pagesIncludingCount"`
			AvailableSizesCount int `json:"availableSizesCount"`
		} `json:"insightsMetadata"`
		ImageID     string `json:"imageId"`
		AccentColor string `json:"accentColor"`
	} `json:"value"`
	PivotSuggestions []struct {
		Pivot       string        `json:"pivot"`
		Suggestions []interface{} `json:"suggestions"`
	} `json:"pivotSuggestions"`
	RelatedSearches []struct {
		Text         string `json:"text"`
		DisplayText  string `json:"displayText"`
		WebSearchURL string `json:"webSearchUrl"`
		SearchLink   string `json:"searchLink"`
		Thumbnail    struct {
			ThumbnailURL string `json:"thumbnailUrl"`
		} `json:"thumbnail"`
	} `json:"relatedSearches"`
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
	queryUrl := baseUrl + url.QueryEscape(query)
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
		return "", errors.New("no images found")
	}
	return data.Value[0].ThumbnailURL, nil
}
