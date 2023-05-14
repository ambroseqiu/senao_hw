package repository

import (
	"context"
	"testing"

	"github.com/ambroseqiu/senao_hw/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func getRandomUser(t *testing.T) *User {
	userName := util.RandomString(6)
	password := util.RandomPassword(8)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)
	return &User{
		Username:       userName,
		HashedPassword: hashedPassword,
	}
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)

	user := getRandomUser(t)

	mockRepo.EXPECT().CreateUser(gomock.Any(), user).Return(nil)

	err := mockRepo.CreateUser(context.Background(), user)
	require.NoError(t, err)
}

func TestCreateUserAlreadyExisted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepository(ctrl)

	user := getRandomUser(t)

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), user).
		Return(nil)

	err := mockRepo.CreateUser(context.Background(), user)
	require.NoError(t, err)

	mockRepo.EXPECT().
		CreateUser(gomock.Any(), user).
		Return(ErrUserIsAlreadyExisted)

	err = mockRepo.CreateUser(context.Background(), user)
	require.EqualError(t, err, ErrUserIsAlreadyExisted.Error())
}
