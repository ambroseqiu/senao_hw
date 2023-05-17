package main

import (
	"fmt"
	"os"

	"github.com/ambroseqiu/senao_hw/controller"
	"github.com/ambroseqiu/senao_hw/docs"
	"github.com/ambroseqiu/senao_hw/migrations"
	"github.com/ambroseqiu/senao_hw/model"
	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "0.0.0.0:8080"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	gormDB, err := repository.NewGormDB()
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	migrations.RunMigration(gormDB)

	repo := repository.NewAccountRepository(gormDB)
	usecase := model.NewUsecaseHandler(repo)
	controller := controller.NewController(usecase)
	route := gin.Default()
	controller.SetRoute(route)
	httpHost := fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("HTTP_PORT"))
	controller.Start(httpHost)
}
