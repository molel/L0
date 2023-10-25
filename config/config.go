package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	HttpPort         int    `json:"http_port"`
	PostgresPort     int    `json:"postgres_port"`
	PostgresUser     string `json:"postgres_user"`
	PostgresPassword string `json:"postgres_password"`
	PostgresDb       string `json:"postgres_db"`
	StanClusterId    string `json:"stan_cluster_id"`
}

func NewConfig() Config {
	file, err := os.ReadFile(".\\..\\config\\configs.json")
	if err != nil {
		log.Fatalf("Cannnot read config file: %s", err)
	}

	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Cannnot parse config file: %s", err)
	}

	return config
}
