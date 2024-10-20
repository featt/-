package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	phone string
	apiHash string
	apiId int
}


func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	phone := os.Getenv("BOT_TOKEN")
	if phone == "" {
		return nil, errors.New("no phone")
	}
	apiId, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		return nil, err
	}
	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		return nil, err
	}
	return &Config{
		phone: phone,
		apiHash: apiHash,
		apiId: apiId,
	}, nil
}

func (c Config) Phone() string {
	return c.phone
}

func (c Config) ApiHash() string {
	return c.apiHash
}

func (c Config) ApiId() int {
	return c.apiId
}
