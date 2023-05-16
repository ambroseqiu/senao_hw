package model

import (
	"context"
	"testing"

	"github.com/ambroseqiu/senao_hw/repository"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	correctUsername := util.RandomString(8)
	correctPassword := util.RandomPassword(8)

	testCase := []struct {
		name             string
		request          AccountRequest
		setMockExpection func(mockRepo *repository.MockAccountRepository)
		verify           func(rsp *AccountResponse, err error)
	}{
		{
			name: "ok",
			request: AccountRequest{
				Username: correctUsername,
				Password: correctPassword,
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.NoError(t, err)
				require.True(t, rsp.Success)
				require.Equal(t, "", rsp.Reason)
			},
		},
		{
			name: "username is too short",
			request: AccountRequest{
				Username: util.RandomString(1),
				Password: correctPassword,
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.EqualError(t, err, ErrAccountRequestValidationFailed.Error())
				require.False(t, rsp.Success)
				require.Equal(t, ErrUsernameIsTooShort.Error(), rsp.Reason)
			},
		},
		{
			name: "password is too short",
			request: AccountRequest{
				Username: correctUsername,
				Password: util.RandomPassword(6),
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.EqualError(t, err, ErrAccountRequestValidationFailed.Error())
				require.False(t, rsp.Success)
				require.Equal(t, ErrPasswordIsTooShort.Error(), rsp.Reason)
			},
		},
		{
			name: "wrong password format",
			request: AccountRequest{
				Username: correctUsername,
				Password: "12345678",
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.EqualError(t, err, ErrAccountRequestValidationFailed.Error())
				require.False(t, rsp.Success)
				require.Equal(t, ErrPasswordValidationFailed.Error(), rsp.Reason)
			},
		},
		{
			name: "account is duplicated",
			request: AccountRequest{
				Username: correctUsername,
				Password: correctPassword,
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(repository.ErrAccountIsDuplicated)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.EqualError(t, err, ErrAccountIsAlreadyExisted.Error())
				require.False(t, rsp.Success)
				require.Equal(t, ErrAccountIsAlreadyExisted.Error(), rsp.Reason)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockRepo := repository.NewMockAccountRepository(ctrl)
			usecase := NewUsecaseHandler(mockRepo)

			tc.setMockExpection(mockRepo)

			rsp, err := usecase.CreateAccount(context.Background(), tc.request)
			tc.verify(rsp, err)
		})
	}

}

func TestLoginAccount(t *testing.T) {
	correctUsername := util.RandomString(8)
	correctPassword := util.RandomPassword(8)
	hashedPassword, err := util.HashedPassword(correctPassword)
	require.NoError(t, err)

	testCase := []struct {
		name             string
		request          AccountRequest
		setMockExpection func(mockRepo *repository.MockAccountRepository)
		verify           func(rsp *AccountResponse, err error)
	}{
		{
			name: "ok",
			request: AccountRequest{
				Username: correctUsername,
				Password: correctPassword,
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().GetAccount(gomock.Any(), correctUsername).
					Return(&repository.Account{
						Username:       correctUsername,
						HashedPassword: hashedPassword,
					}, nil)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.NoError(t, err)
				require.True(t, rsp.Success)
				require.Equal(t, "", rsp.Reason)
			},
		},
		{
			name: "account record not found",
			request: AccountRequest{
				Username: correctUsername,
				Password: correctPassword,
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().GetAccount(gomock.Any(), correctUsername).
					Return(nil, repository.ErrAccountRecordNotFound)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.EqualError(t, err, ErrLoginAccountNotFound.Error())
				require.False(t, rsp.Success)
				require.Equal(t, ErrLoginAccountNotFound.Error(), rsp.Reason)
			},
		},
		{
			name: "wrong password",
			request: AccountRequest{
				Username: correctUsername,
				Password: util.RandomPassword(8),
			},
			setMockExpection: func(mockRepo *repository.MockAccountRepository) {
				mockRepo.EXPECT().GetAccount(gomock.Any(), correctUsername).
					Return(&repository.Account{
						Username:       correctUsername,
						HashedPassword: hashedPassword,
					}, nil)
			},
			verify: func(rsp *AccountResponse, err error) {
				require.EqualError(t, err, ErrLoginWrongPassword.Error())
				require.False(t, rsp.Success)
				require.Equal(t, ErrLoginWrongPassword.Error(), rsp.Reason)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockRepo := repository.NewMockAccountRepository(ctrl)
			usecase := NewUsecaseHandler(mockRepo)

			tc.setMockExpection(mockRepo)

			rsp, err := usecase.LoginAccount(context.Background(), tc.request)
			tc.verify(rsp, err)
		})
	}
}

func TestLoginFailedAttemptFiveTime(t *testing.T) {
	username := util.RandomString(10)
	password := util.RandomPassword(10)
	wrongPassword := util.RandomPassword(10)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)

	mockRepo := repository.NewMockAccountRepository(ctrl)

	usecase := NewUsecaseHandler(mockRepo)

	mockRepo.EXPECT().GetAccount(gomock.Any(), username).Times(5).Return(&repository.Account{
		Username:       username,
		HashedPassword: hashedPassword,
	}, nil)

	req := AccountRequest{
		Username: username,
		Password: wrongPassword,
	}

	for i := 1; i <= 5; i++ {
		rsp, err := usecase.LoginAccount(context.Background(), req)
		require.EqualError(t, err, ErrLoginWrongPassword.Error())
		require.False(t, rsp.Success)
		require.Equal(t, ErrLoginWrongPassword.Error(), rsp.Reason)
	}

	rsp, err := usecase.LoginAccount(context.Background(), req)
	require.EqualError(t, err, ErrLoginAttemptBlocked.Error())
	require.False(t, rsp.Success)
	require.Equal(t, ErrLoginAttemptBlocked.Error(), rsp.Reason)
}
