package images

import (
	"github.com/psidex/PortsmouthShippingMovements/internal/bing"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// ShipImageUrlStorage manages the storage and caching of ship image urls.
// Urls are stored in memory and on disk as files, this ensures speed and persistence.
type ShipImageUrlStorage struct {
	imageSearchApi bing.ImageSearchApi // The API for searching for images.
	path           string              // Where the images URLs are to be stored (as files).
	memory         map[string]string   // An in-memory map of known ship image URLs.
}

// NewShipImageUrlStorage creates a new ShipImageUrlStorage, can return an error if there is a problem with imageDirectoryPath.
func NewShipImageUrlStorage(imageSearchApi bing.ImageSearchApi, storagePath string) (ShipImageUrlStorage, error) {
	_, err := os.Stat(storagePath)
	if os.IsNotExist(err) {
		if err = os.Mkdir(storagePath, 0644); err != nil {
			return ShipImageUrlStorage{}, err
		}
		log.Printf("Created url storage directory: %s", storagePath)
	} else if err != nil {
		return ShipImageUrlStorage{}, err
	}

	return ShipImageUrlStorage{
		imageSearchApi: imageSearchApi,
		path:           storagePath,
		memory:         make(map[string]string),
	}, nil
}

// saveUrlToFile read a url from a file with the ship name as the file name.
func (s ShipImageUrlStorage) readUrlFromFile(shipName string) (string, error) {
	dat, err := ioutil.ReadFile(path.Join(s.path, shipName))
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// saveUrlToFile writes the given url to a file with the ship name as the file name.
func (s ShipImageUrlStorage) saveUrlToFile(shipName, imageUrl string) error {
	data := []byte(imageUrl)
	err := ioutil.WriteFile(path.Join(s.path, shipName), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetUrlForShip takes a ship name and returns a string containing the URL to an image of that ship.
// If an error occurs or an image can't be found, an empty string ("") is returned.
func (s ShipImageUrlStorage) GetUrlForShip(shipName string) string {
	if url, ok := s.memory[shipName]; ok {
		return url
	}

	if url, err := s.readUrlFromFile(shipName); err == nil {
		s.memory[shipName] = url
		return url
	}

	log.Printf("Searching BingApi API for images of ship: %s", shipName)

	url, err := s.imageSearchApi.SearchForImage(shipName)
	if err != nil {
		log.Printf("Error searching for image: %v", err)
		return ""
	}

	if url != "" {
		// Only save if we actually have a url.
		err := s.saveUrlToFile(shipName, url)
		if err == nil {
			// If there was an error writing to the file, not saving the url in mem will trigger another write attempt
			// next time.
			s.memory[shipName] = url
		}
	}

	return url
}
