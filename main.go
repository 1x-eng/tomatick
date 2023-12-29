package main

import (
	"log"

	"github.com/1x-eng/tomatick/config"

	"github.com/1x-eng/tomatick/cmd"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cmd.Execute(cfg)
}
