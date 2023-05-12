package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *apiController) GetHandler(ctx *gin.Context) {

	if err := ctrl.usecase.GetApi(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, "success")
}
