package main

import (
	"fmt"
	"os"

	"github.com/ambroseqiu/senao_hw/controller"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	controller := controller.NewController()
	controller.SetRoute()
	httpHost := fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("HTTP_PORT"))
	controller.Start(httpHost)
}
