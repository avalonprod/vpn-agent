package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WGInterface     string
	BaseIP          string
	Endpoint        string
	AuthToken       string
	DNSServer       string
	ServerPublicKey string
	Port            string
}

var ServerConfig Config

func mustLoad() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load .env file. err: %v", err)
	}

	ServerConfig = Config{
		WGInterface:     os.Getenv("WG_INTERFACE"),
		BaseIP:          os.Getenv("BASE_IP"),
		Endpoint:        os.Getenv("ENDPOINT"),
		AuthToken:       os.Getenv("AUTH_TOKEN"),
		DNSServer:       os.Getenv("DNS_SERVER"),
		ServerPublicKey: os.Getenv("SERVER_PUBLIC_KEY"),
		Port:            os.Getenv("PORT"),
	}
}
