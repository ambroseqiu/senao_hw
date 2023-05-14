package main

import (
	"fmt"
	"os"

	"github.com/ambroseqiu/senao_hw/controller"
	"github.com/ambroseqiu/senao_hw/migrations"
	"github.com/ambroseqiu/senao_hw/model"
	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	gormDB, err := repository.NewGormDB()
	if err != nil {
		log.Fatal().Err(err)
	}

	migrations.RunMigration(gormDB)

	repo := repository.NewAccountRepository(gormDB)
	usecase := model.NewUsecaseHandler(repo)
	controller := controller.NewController(usecase)
	route := gin.Default()
	controller.SetRoute(route)
	httpHost := fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("HTTP_PORT"))
	controller.Start(httpHost)
}
