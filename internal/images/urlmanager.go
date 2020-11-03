package images

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/bing"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// UrlManager manages the search and caching of image urls.
// Urls are stored in memory and on disk as files, this ensures speed and a persistent cache.
// We store and use image URLs instead of the actual images so we don't have to worry about hosting them (licensing, bandwidth issues, etc.).
type UrlManager struct {
	imageSearchApi bing.ImageSearchApi // The API for searching for images.
	path           string              // Where the images URLs are to be stored (as files).
	memory         map[string]string   // An in-memory map of known image URLs.
}

// NewUrlManager creates a new UrlManager, can return an error if there is a problem with imageDirectoryPath.
func NewUrlManager(imageSearchApi bing.ImageSearchApi, storagePath string) (UrlManager, error) {
	_, err := os.Stat(storagePath)
	if os.IsNotExist(err) {
		if err = os.Mkdir(storagePath, 0644); err != nil {
			return UrlManager{}, err
		}
		log.Printf("Created image url storage directory: %s", storagePath)
	} else if err != nil {
		return UrlManager{}, err
	}

	return UrlManager{
		imageSearchApi: imageSearchApi,
		path:           storagePath,
		memory:         make(map[string]string),
	}, nil
}

// saveUrlToFile read a url from a file with the query as the file name.
func (m UrlManager) readUrlFromFile(query string) (string, error) {
	dat, err := ioutil.ReadFile(path.Join(m.path, query))
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// saveUrlToFile writes the given url to a file with the query as the file name.
func (m UrlManager) saveUrlToFile(query, imageUrl string) error {
	data := []byte(imageUrl)
	err := ioutil.WriteFile(path.Join(m.path, query), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetUrl returns a string containing the URL to a thumbnail image for the given query.
// If an error occurs or an image can't be found, an empty string ("") is returned.
func (m UrlManager) GetUrl(query string) string {
	if url, ok := m.memory[query]; ok {
		return url
	}

	if url, err := m.readUrlFromFile(query); err == nil {
		m.memory[query] = url
		return url
	}

	log.Printf("Searching Bing API for images of: %s", query)
	url, err := m.imageSearchApi.SearchForImage(query)
	if err != nil {
		log.Printf("Error searching for image: %s", err)
		return ""
	}

	err = m.saveUrlToFile(query, url)
	if err != nil {
		log.Printf("Error writing image url to file: %s", err)
		// Don't save the url in memory to trigger another write attempt next time.
		return url
	}

	m.memory[query] = url
	return url
}
