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

// CreateAccount godoc
// @Summary      Create an account
// @Description  Create account by username and password
// @Tags         accounts
// @Param        accountRequest body model.AccountRequest true "Account Request Struct"
// @Success      200  {object}  model.AccountResponse
// @Failure      400  {object}  model.AccountResponse
// @Failure      500  {object}  error
// @Router       /accounts [post]
func (ctrl *apiController) CreateAccount(ctx *gin.Context) {
	var req model.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := ctrl.usecase.CreateAccount(context.Background(), req)
	if err != nil {
		if err == model.ErrAccountRequestValidationFailed {
			ctx.JSON(http.StatusBadRequest, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

// LoginAccount godoc
// @Summary      Login account
// @Description  login account and verify username and password
// @Tags         accounts
// @Param        accountRequest body model.AccountRequest true "Account Request Struct"
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.AccountResponse
// @Router       /login [post]
func (ctrl *apiController) LoginAccount(ctx *gin.Context) {
	var req model.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := ctrl.usecase.LoginAccount(context.Background(), req)
	if err != nil {
		if err == model.ErrAccountRequestValidationFailed || err == model.ErrLoginAccountNotFound {
			ctx.JSON(http.StatusBadRequest, rsp)
		} else if err == model.ErrLoginAccountNotAllowed {
			ctx.JSON(http.StatusUnauthorized, rsp)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
