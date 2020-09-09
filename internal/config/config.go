package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ImageStoragePath      string `mapstructure:"image_storage_path"`
	AccessLogPath         string `mapstructure:"access_log_path"`
	BingImageSearchApiKey string `mapstructure:"bing_image_search_api_key"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("image_storage_path", "./imagestorage")
	viper.SetDefault("access_log_path", "./access.log")
	viper.SetDefault("bing_image_search_api_key", "")

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
