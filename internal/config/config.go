package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
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
		viper.SetDefault("bing_image_search_api_key", "")
		log.Println("Running inside Docker")
	} else {
		viper.SetDefault("image_storage_path", "./imagestorage")
		viper.SetDefault("access_log_path", "./access.log")
		viper.SetDefault("bing_image_search_api_key", "")
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

	return config, nil
}
