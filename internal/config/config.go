package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	ContactEmail          string `mapstructure:"contact_email"`
	ImageStoragePath      string `mapstructure:"image_storage_path"`
	AccessLogPath         string `mapstructure:"access_log_path"`
	BingImageSearchApiKey string `mapstructure:"bing_image_search_api_key"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	insideDocker := os.Getenv("INSIDE_DOCKER")
	if insideDocker != "" {
		// If running in Docker and you want to mount a volume, mount it at "/data".
		// If doing this, you probably also want to prepend the other storage paths with "/data/".
		viper.AddConfigPath("/data")
		viper.SetDefault("image_storage_path", "/data/imagestorage")
		viper.SetDefault("access_log_path", "/data/access.log")
		log.Println("Running inside Docker")
	} else {
		viper.SetDefault("image_storage_path", "./imagestorage")
		viper.SetDefault("access_log_path", "./access.log")
		log.Println("Not running inside Docker")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	log.Println("Configuration loaded")
	log.Printf("contact_email: %s", config.ContactEmail)
	log.Printf("image_storage_path: %s", config.ImageStoragePath)
	log.Printf("access_log_path: %s", config.AccessLogPath)

	return config, nil
}
