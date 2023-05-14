package repository

import (
	"context"
	"testing"

	"github.com/ambroseqiu/senao_hw/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getRandomAccount(t *testing.T) *Account {
	userName := util.RandomString(6)
	password := util.RandomPassword(8)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)
	return &Account{
		Username:       userName,
		HashedPassword: hashedPassword,
	}
}

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockAccountRepository(ctrl)

	account := getRandomAccount(t)

	mockRepo.EXPECT().CreateAccount(gomock.Any(), account).Return(nil)

	err := mockRepo.CreateAccount(context.Background(), account)
	require.NoError(t, err)
}

func TestCreateAccountAlreadyExisted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockAccountRepository(ctrl)

	account := getRandomAccount(t)

	mockRepo.EXPECT().
		CreateAccount(gomock.Any(), account).
		Return(nil)

	err := mockRepo.CreateAccount(context.Background(), account)
	require.NoError(t, err)

	mockRepo.EXPECT().
		CreateAccount(gomock.Any(), account).
		Return(ErrAccountIsAlreadyExisted)

	err = mockRepo.CreateAccount(context.Background(), account)
	require.EqualError(t, err, ErrAccountIsAlreadyExisted.Error())
}
