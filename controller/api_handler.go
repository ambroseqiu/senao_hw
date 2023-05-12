package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *apiController) GetHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "success")
}
