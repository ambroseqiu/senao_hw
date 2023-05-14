package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ambroseqiu/senao_hw/model"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreatAccount(t *testing.T) {
	username := util.RandomString(10)
	password := util.RandomPassword(8)

	testCase := []struct {
		name             string
		body             gin.H
		setMockExpection func(mockUsecase *model.MockUsecaseHandler)
		checkResponse    func(*httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username": username,
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Return(&model.AccountResponse{
						Success: true,
						Reason:  "",
					}, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
				byteData, err := ioutil.ReadAll(rr.Body)
				require.NoError(t, err)

				rsp := &model.AccountResponse{}
				err = json.Unmarshal(byteData, rsp)
				require.NoError(t, err)

				require.True(t, rsp.Success)
			},
		},
		{
			name: "bad request",
			body: gin.H{
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "request validation failed",
			body: gin.H{
				"username": util.RandomString(2),
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Return(&model.AccountResponse{
						Success: false,
						Reason:  model.ErrAccountRequestValidationFailed.Error(),
					}, model.ErrAccountRequestValidationFailed)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
				byteData, err := ioutil.ReadAll(rr.Body)
				require.NoError(t, err)

				rsp := &model.AccountResponse{}
				err = json.Unmarshal(byteData, rsp)
				require.NoError(t, err)
				require.False(t, rsp.Success)
				require.Equal(t, model.ErrAccountRequestValidationFailed.Error(), rsp.Reason)
			},
		},
		{
			name: "internal server error",
			body: gin.H{
				"username": username,
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("unknown error"))
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUsecase := model.NewMockUsecaseHandler(ctrl)
			controller := NewController(mockUsecase)
			route := gin.Default()
			controller.SetRoute(route)
			url := "/api/account"

			tc.setMockExpection(mockUsecase)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(data))
			r := httptest.NewRecorder()
			route.ServeHTTP(r, httpReq)

			tc.checkResponse(r)
		})
	}
}

func TestLoginAccount(t *testing.T) {
	username := util.RandomString(10)
	password := util.RandomPassword(8)

	testCase := []struct {
		name             string
		body             gin.H
		setMockExpection func(mockUsecase *model.MockUsecaseHandler)
		checkResponse    func(*httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username": username,
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().LoginAccount(gomock.Any(), gomock.Any()).
					Return(&model.AccountResponse{
						Success: true,
						Reason:  "",
					}, nil)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
				byteData, err := ioutil.ReadAll(rr.Body)
				require.NoError(t, err)

				rsp := &model.AccountResponse{}
				err = json.Unmarshal(byteData, rsp)
				require.NoError(t, err)

				require.True(t, rsp.Success)
			},
		},
		{
			name: "bad request",
			body: gin.H{},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().LoginAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
		{
			name: "username is not existed",
			body: gin.H{
				"username": util.RandomString(8),
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().LoginAccount(gomock.Any(), gomock.Any()).
					Return(&model.AccountResponse{
						Success: false,
						Reason:  model.ErrLoginAccountNotFound.Error(),
					}, model.ErrLoginAccountNotFound)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
				byteData, err := ioutil.ReadAll(rr.Body)
				require.NoError(t, err)

				rsp := &model.AccountResponse{}
				err = json.Unmarshal(byteData, rsp)
				require.NoError(t, err)

				require.False(t, rsp.Success)
				require.Equal(t, model.ErrLoginAccountNotFound.Error(), rsp.Reason)
			},
		},
		{
			name: "login is not allowed",
			body: gin.H{
				"username": username,
				"password": "wrong password",
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().LoginAccount(gomock.Any(), gomock.Any()).
					Return(&model.AccountResponse{
						Success: false,
						Reason:  model.ErrLoginAccountNotAllowed.Error(),
					}, model.ErrLoginAccountNotAllowed)
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rr.Code)
				byteData, err := ioutil.ReadAll(rr.Body)
				require.NoError(t, err)

				rsp := &model.AccountResponse{}
				err = json.Unmarshal(byteData, rsp)
				require.NoError(t, err)

				require.False(t, rsp.Success)
				require.Equal(t, model.ErrLoginAccountNotAllowed.Error(), rsp.Reason)
			},
		},
		{
			name: "internal server error",
			body: gin.H{
				"username": username,
				"password": password,
			},
			setMockExpection: func(mockUsecase *model.MockUsecaseHandler) {
				mockUsecase.EXPECT().LoginAccount(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("unknown error"))
			},
			checkResponse: func(rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUsecase := model.NewMockUsecaseHandler(ctrl)
			controller := NewController(mockUsecase)
			route := gin.Default()
			controller.SetRoute(route)
			url := "/api/login"

			tc.setMockExpection(mockUsecase)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			httpReq, _ := http.NewRequest("GET", url, bytes.NewReader(data))
			r := httptest.NewRecorder()
			route.ServeHTTP(r, httpReq)

			tc.checkResponse(r)
		})
	}

}
