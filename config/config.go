package config

import (
    "log"

    "github.com/spf13/viper"
)

type Config struct {
    Port   int
    DBUrl  string
    Email  string
    Pass   string
    SID    string
    Token  string
    Phone  string
    URL    string
    JwtKey string
}

func LoadConfig() *Config {
    // Viper ko bata rahe hain kaha dhoondhna hai
    viper.AddConfigPath(".")        // project root
    viper.AddConfigPath("./config") // config folder
    viper.SetConfigName("config")   // config.json
    viper.SetConfigType("json")

    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("❌ Error loading config file: %v", err)
    }

    log.Printf("DEBUG: Loaded DB_URL: '%s'", viper.GetString("DB_URL"))

    return &Config{
        Port:   viper.GetInt("PORT"),
        DBUrl:  viper.GetString("DB_URL"),
        Email:  viper.GetString("EMAIL"),
        Pass:   viper.GetString("PASS"),
        SID:    viper.GetString("SID"),
        Token:  viper.GetString("TOKEN"),
        Phone:  viper.GetString("PHONE"),
        URL:    viper.GetString("URL"),
        JwtKey: viper.GetString("JWT_KEY"),
    }
}
