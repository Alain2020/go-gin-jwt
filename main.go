package main

import (
	"errors"
	"go-gin-jwt/app"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//Check if .env file is exist or not
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		godotenv.Load()
	}
}

func main() {
	app := app.NewApp()
	app.Run()
}
