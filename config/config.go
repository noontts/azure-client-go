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

type Config struct {
	Azure  AzureConfig
	Client ClientConfig
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
	}

	missing := []string{}
	// if cfg.Azure.StorageAccount == "" {
	// 	missing = append(missing, "AZURE_STORAGE_ACCOUNT")
	// }
	// if cfg.Azure.EmailEndpoint == "" {
	// 	missing = append(missing, "AZURE_EMAIL_ENDPOINT")
	// }
	// if cfg.Azure.EmailAccessKey == "" {
	// 	missing = append(missing, "AZURE_EMAIL_ACCESS_KEY")
	// }
	if cfg.Client.PaymentBaseURL == "" {
		missing = append(missing, "PAYMENT_BASE_URL")
	}
	if cfg.Client.MemberBaseURL == "" {
		missing = append(missing, "MEMBER_BASE_URL")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}
	return cfg, nil
}
