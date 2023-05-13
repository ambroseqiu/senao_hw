package controller

import (
	"context"
	"net/http"

	"github.com/ambroseqiu/senao_hw/model"
	"github.com/gin-gonic/gin"
)

func errResponse(err error) *gin.H {
	return &gin.H{"err": err}
}

func (ctrl *apiController) CreateUser(ctx *gin.Context) {
	var req model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	rsp, err := ctrl.usecase.CreateUser(context.Background(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
