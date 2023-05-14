package controller

import (
	"context"
	"net/http"

	"github.com/ambroseqiu/senao_hw/model"
	"github.com/gin-gonic/gin"
)

func errResponse(err error) *gin.H {
	return &gin.H{"err": err.Error()}
}

func (ctrl *apiController) CreateUser(ctx *gin.Context) {
	var req model.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := ctrl.usecase.CreateUser(context.Background(), req)
	if err != nil {
		if err == model.ErrCreateUserRequestValidationFailed {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (ctrl *apiController) LoginUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
