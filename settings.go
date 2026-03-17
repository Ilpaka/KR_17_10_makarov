package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Settings struct {
	JwtSecret string
}

var settings *Settings

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	settings = &Settings{
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
}
