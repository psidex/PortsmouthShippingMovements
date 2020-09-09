package images

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

// ShipImageUrlStorage manages the storage and caching of ship image urls.
// Urls are stored in memory and on disk as files, this ensures speed and persistence.
type ShipImageUrlStorage struct {
	bingApiKey string            // The Bing API key for searching for images.
	path       string            // Where the images URLs are to be stored (as files).
	memory     map[string]string // An in-memory map of known ship image URLs.
}

// NewShipImageUrlStorage creates a new ShipImageUrlStorage, can return an error if there is a problem with imageDirectoryPath.
func NewShipImageUrlStorage(bingApiKey, imageUrlDirectoryPath string) (ShipImageUrlStorage, error) {
	_, err := os.Stat(imageUrlDirectoryPath)
	if os.IsNotExist(err) {
		if err = os.Mkdir(imageUrlDirectoryPath, 0644); err != nil {
			return ShipImageUrlStorage{}, err
		}
		log.Printf("Created url storage directory: %s", imageUrlDirectoryPath)
	} else if err != nil {
		return ShipImageUrlStorage{}, err
	}

	return ShipImageUrlStorage{
		bingApiKey: bingApiKey,
		path:       imageUrlDirectoryPath,
		memory:     make(map[string]string),
	}, nil
}

// saveUrlToFile read a url from a file with the ship name as the file name.
func (i ShipImageUrlStorage) readUrlFromFile(shipName string) (string, error) {
	dat, err := ioutil.ReadFile(path.Join(i.path, shipName))
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// saveUrlToFile writes the given url to a file with the ship name as the file name.
func (i ShipImageUrlStorage) saveUrlToFile(shipName, imageUrl string) error {
	data := []byte(imageUrl)
	err := ioutil.WriteFile(path.Join(i.path, shipName), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetUrlForShip takes a ship name and returns a string containing the URL to an image of that ship.
func (i ShipImageUrlStorage) GetUrlForShip(shipName string) string {
	if url, ok := i.memory[shipName]; ok {
		return url
	}

	// If the file doesn't exist or there is a problem, move on from this block.
	if url, err := i.readUrlFromFile(shipName); url != "" && err == nil {
		i.memory[shipName] = url
		return url
	}

	log.Printf("Searching Bing API for images of ship: %s", shipName)
	url := searchForShipImage(i.bingApiKey, shipName)
	if url != "" {
		// We will have to return an empty string is it is one, but only save if it's not one.
		err := i.saveUrlToFile(shipName, url)
		if err == nil {
			// If there was an error writing to the file, not saving the url in mem will trigger another write attempt
			// next time.
			i.memory[shipName] = url
		}
	}
	return url
}
