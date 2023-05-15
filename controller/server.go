package controller

import (
	"github.com/ambroseqiu/senao_hw/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type apiController struct {
	usecase model.UsecaseHandler
	route   *gin.Engine
}

func NewController(usecase model.UsecaseHandler) apiController {
	return apiController{
		usecase: usecase,
	}
}

func (ctrl *apiController) SetRoute(route *gin.Engine) {

	// 設置CORS中間件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token"}
	route.Use(cors.New(config))

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiRoute := route.Group("/api")
	apiRoute.POST("/accounts", ctrl.CreateAccount)
	apiRoute.POST("/login", ctrl.LoginAccount)

	ctrl.route = route
}

func (ctrl *apiController) Start(addr string) error {
	return ctrl.route.Run(addr)
}
