package controller

import (
	"github.com/ambroseqiu/senao_hw/model"
	"github.com/gin-gonic/gin"
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

func (ctrl *apiController) SetRoute() {
	route := gin.Default()

	apiRoute := route.Group("/api")
	apiRoute.POST("/user", ctrl.CreateUser)

	ctrl.route = route
}

func (ctrl *apiController) Start(addr string) error {
	return ctrl.route.Run(addr)
}
