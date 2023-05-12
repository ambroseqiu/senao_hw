package controller

import "github.com/gin-gonic/gin"

type apiController struct {
	route *gin.Engine
}

func NewController() apiController {
	return apiController{}
}

func (ctrl *apiController) SetRoute() {
	route := gin.Default()

	route.GET("/", ctrl.GetHandler)

	ctrl.route = route
}

func (ctrl *apiController) Start(addr string) error {
	return ctrl.route.Run(addr)
}
