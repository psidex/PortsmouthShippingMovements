package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config contains the application configuration, to be unmarshalled into by Viper.
type Config struct {
	UpdateCronString      string `mapstructure:"update_cron_string"`
	ContactEmail          string `mapstructure:"contact_email"`
	StoragePath           string `mapstructure:"storage_path"`
	BingImageSearchApiKey string `mapstructure:"bing_image_search_api_key"`
}

// Get looks for and read any config file found into the Config struct.
func Get() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Update stored movements on the hour at midnight, 6am, midday, and 6pm.
	viper.SetDefault("update_cron_string", "0 0,6,12,18 * * *")
	viper.SetDefault("storage_path", "./storage")

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
	log.Printf("update_cron_string: %s", config.UpdateCronString)
	log.Printf("contact_email: %s", config.ContactEmail)
	log.Printf("storage_path: %s", config.StoragePath)

	return config, nil
}
