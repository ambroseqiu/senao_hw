package controller

import (
	"context"
	"net/http"

	_ "github.com/ambroseqiu/senao_hw/docs"
	"github.com/ambroseqiu/senao_hw/model"
	"github.com/gin-gonic/gin"
)

func errResponse(err error) *gin.H {
	return &gin.H{"err": err.Error()}
}

// CreateAccount godoc
// @Summary      Create an account
// @Description  Create account by username and password
// @Description  Note:
// @Description  username: a string representing the desired username for the account, with a minimum length of 3 characters and a maximum length of 32 characters.
// @Description  password: a string representing the desired password for the account, with a minimum length of 8 characters and a maximum length of 32 characters,
// @Description  containing at least 1 uppercase letter, 1 lowercase letter, and 1 number.
// @Tags         accounts
// @Param        accountRequest body model.AccountRequest true "Account Request Struct"
// @Success      200  {object}  model.DocResponseSuccess
// @Failure      400  {object}  model.DocResponseBadRequest
// @Failure      409  {object}  model.DocResponseAlreadyExisted "Account Is Already Existed"
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
		} else if err == model.ErrAccountIsAlreadyExisted {
			ctx.JSON(http.StatusConflict, rsp)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

// LoginAccount godoc
// @Summary      Login account
// @Description  Login account and verify username and password
// @Description  Note:
// @Description  If the password verification fails five times, the user should wait one minute before attempting to verify the password again.
// @Tags         accounts
// @Param        accountRequest body model.AccountRequest true "Account Request Struct"
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.DocResponseSuccess
// @Failure      400  {object}  model.DocResponseAccountNotFound
// @Failure      401  {object}  model.DocResponseWrongPassword
// @Failure      429  {object}  model.DocResponseTooManyRequest "Too Many Failed Login Attempts"
// @Router       /login [post]
func (ctrl *apiController) LoginAccount(ctx *gin.Context) {
	var req model.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	rsp, err := ctrl.usecase.LoginAccount(context.Background(), req)
	if err != nil {
		if err == model.ErrLoginAccountNotFound {
			ctx.JSON(http.StatusBadRequest, rsp)
		} else if err == model.ErrLoginWrongPassword {
			ctx.JSON(http.StatusUnauthorized, rsp)
		} else if err == model.ErrLoginAttemptBlocked {
			ctx.JSON(http.StatusTooManyRequests, rsp)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
