package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config struct to store loaded configuration values
type Config struct {
	Port                int
	DBUrl               string
	Email               string
	Pass                string
	SID                 string
	Token               string
	Phone               string
	URL                 string
	JwtKey              string
	CloudinaryCloudName string
	CloudinaryApiKey    string
	CloudinaryApiSecret string
}

// LoadConfig loads the configuration file using Viper
func LoadConfig() *Config {
	viper.SetConfigName("config")  // File name without extension
	viper.SetConfigType("json")    // Config file type
	viper.AddConfigPath("config")  // Path where config file is stored

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	// Map values from Viper to our Config struct
	return &Config{
		Port:                viper.GetInt("PORT"),
		DBUrl:               viper.GetString("DB_URL"),
		Email:               viper.GetString("EMAIL"),
		Pass:                viper.GetString("PASS"),
		SID:                 viper.GetString("SID"),
		Token:               viper.GetString("TOKEN"),
		Phone:               viper.GetString("PHONE"),
		URL:                 viper.GetString("URL"),
		JwtKey:              viper.GetString("JWT_KEY"),
		CloudinaryCloudName: viper.GetString("CLOUDINARY_CLOUD_NAME"),
		CloudinaryApiKey:    viper.GetString("CLOUDINARY_API_KEY"),
		CloudinaryApiSecret: viper.GetString("CLOUDINARY_API_SECRET"),
	}
}