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

func (ctrl *apiController) CreateAccount(ctx *gin.Context) {
	var req model.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := ctrl.usecase.CreateAccount(context.Background(), req)
	if err != nil {
		if err == model.ErrCreateAccountRequestValidationFailed {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (ctrl *apiController) LoginAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}
