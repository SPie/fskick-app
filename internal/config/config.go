package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spie/fskick/internal/db"
)

type AppConfig struct {
    ApiHost string
    DbConfig db.DbConfig
}

func LoadCliConfig() (AppConfig, error) {
    cfg, err := loadEnvConfig()
    if err != nil {
        return AppConfig{}, err
    }

    setDbConfig(&cfg)

    return cfg, nil
}

func LoadApiConfig() (AppConfig, error) {
    cfg, err := loadEnvConfig()
    if err != nil {
        return AppConfig{}, err
    }

    setDbConfig(&cfg)
    setApiConfig(&cfg)

    return cfg, nil
}

func loadEnvConfig() (AppConfig, error) {
    var cfg AppConfig
    err := godotenv.Load()
    if err != nil {
        return AppConfig{}, err
    }

    return cfg, nil
}

func setDbConfig(cfg *AppConfig) {
    cfg.DbConfig = db.CreateDbConfig(
        os.Getenv("DB_DATABASE"),
        os.Getenv("DB_LOG") == "true",
        os.Getenv("DB_DEBUG") == "true",
    )
}

func setApiConfig(cfg *AppConfig) {
    cfg.ApiHost = os.Getenv("API_HOST")
}
