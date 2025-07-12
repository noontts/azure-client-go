package config

import (
	"fmt"
	"os"
)

type AzureConfig struct {
	StorageAccount string
	EmailEndpoint  string
	EmailAccessKey string
}

type ClientConfig struct {
	PaymentBaseURL string
	MemberBaseURL  string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Name     string
}

type Config struct {
	Azure  AzureConfig
	Client ClientConfig
	DB     DBConfig
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Azure: AzureConfig{
			StorageAccount: os.Getenv("AZURE_STORAGE_ACCOUNT"),
			EmailEndpoint:  os.Getenv("AZURE_EMAIL_ENDPOINT"),
			EmailAccessKey: os.Getenv("AZURE_EMAIL_ACCESS_KEY"),
		},
		Client: ClientConfig{
			PaymentBaseURL: os.Getenv("PAYMENT_BASE_URL"),
			MemberBaseURL:  os.Getenv("MEMBER_BASE_URL"),
		},
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("DB_NAME"),
		},
	}

	missing := []string{}
	if cfg.Client.PaymentBaseURL == "" {
		missing = append(missing, "PAYMENT_BASE_URL")
	}
	if cfg.Client.MemberBaseURL == "" {
		missing = append(missing, "MEMBER_BASE_URL")
	}
	if cfg.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DB.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DB.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if cfg.DB.Name == "" {
		missing = append(missing, "DB_NAME")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}
	return cfg, nil
}
