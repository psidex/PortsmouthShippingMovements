package images

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

const apiRoute = "https://api.cognitive.microsoft.com/bing/v7.0/images/search?q="

// imageSearch represents the API response from a Bing image search API request.
// Thanks https://mholt.github.io/json-to-go/
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

// searchForShipImage attempts to find an image URL for the given ship name. Returns "" if none found or an error occurs.
func searchForShipImage(apiKey, shipName string) string {
	client := &http.Client{}
	query := apiRoute + url.QueryEscape(shipName)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		log.Printf("Error creating search request: %v", err)
		return ""
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error requesting search: %v", err)
		return ""
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var data imageSearch
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("Error decoding search response: %v", err)
		return ""
	}

	return data.Value[0].ThumbnailURL
}
