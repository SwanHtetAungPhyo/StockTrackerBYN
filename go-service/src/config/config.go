package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	APIKEY string `json:"apikey"`
	ParameterO Parameter `json:"parameterO"`
}

type Parameter struct{
	BaseUrl string `json:"base_url"`
	Adjusted string `json:"adjusted"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("%v", err.Error())
		return nil, err
	}
	apiKey := os.Getenv("API_KEY")
	return &Config{
		APIKEY: apiKey,
		ParameterO: Parameter{
			BaseUrl:  os.Getenv("BASE_URL"),
			Adjusted: os.Getenv("ADJUSTED"),
		},
	}, nil
}
