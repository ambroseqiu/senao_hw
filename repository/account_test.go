package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ambroseqiu/senao_hw/util"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func setUpAccountMock(t *testing.T) (AccountRepository, *sql.DB, sqlmock.Sqlmock) {
	mockDb, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repo := NewAccountRepository(gormDB)
	return repo, mockDb, mock
}

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
	repo, mockDB, mock := setUpAccountMock(t)
	defer mockDB.Close()

	account := getRandomAccount(t)

	// 模拟插入数据时的返回结果
	rows := sqlmock.NewRows([]string{"id"}).AddRow(account.ID)

	// 设置 mock 预期行为
	mock.ExpectBegin()
	sqlQuery := `INSERT INTO "accounts" ("id","username","hashed_password","created_at","updated_at","deleted_at") 
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	mock.ExpectQuery(sqlQuery).
		WithArgs(account.ID, account.Username, account.HashedPassword, AnyTime{}, AnyTime{}, nil).
		WillReturnRows(rows)
	mock.ExpectCommit()

	err := repo.CreateAccount(context.Background(), account)
	require.NoError(t, err)
}

func TestCreateAccountDuplicate(t *testing.T) {
	repo, mockDB, mock := setUpAccountMock(t)
	defer mockDB.Close()

	account := getRandomAccount(t)

	mock.ExpectBegin()
	sqlQuery := `INSERT INTO "accounts" ("id","username","hashed_password","created_at","updated_at","deleted_at") 
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	mock.ExpectQuery(sqlQuery).
		WithArgs(account.ID, account.Username, account.HashedPassword, AnyTime{}, AnyTime{}, nil).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	err := repo.CreateAccount(context.Background(), account)
	require.EqualError(t, err, ErrAccountIsAlreadyExisted.Error())
}

func TestGetAccount(t *testing.T) {
	repo, mockDB, mock := setUpAccountMock(t)
	defer mockDB.Close()

	account := getRandomAccount(t)

	rows := sqlmock.NewRows([]string{"id", "username", "hashed_password", "created_at", "updated_at", "deleted_at"}).
		AddRow(account.ID, account.Username, account.HashedPassword, time.Now(), time.Now(), nil)

	sqlQuery := `SELECT * FROM "accounts" WHERE username = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`
	mock.ExpectQuery(sqlQuery).
		WithArgs(account.Username).
		WillReturnRows(rows)

	getAccount, err := repo.GetAccount(context.Background(), account.Username)
	require.NoError(t, err)
	require.NotNil(t, getAccount)
}

func TestGetAccountNotExisted(t *testing.T) {
	repo, mockDB, mock := setUpAccountMock(t)
	defer mockDB.Close()

	noExistedName := util.RandomString(3)

	sqlQuery := `SELECT * FROM "accounts" WHERE username = $1 AND "accounts"."deleted_at" IS NULL ORDER BY "accounts"."id" LIMIT 1`
	mock.ExpectQuery(sqlQuery).
		WithArgs(noExistedName).
		WillReturnError(gorm.ErrRecordNotFound)
	getAccount, err := repo.GetAccount(context.Background(), noExistedName)
	require.EqualError(t, err, ErrAccountRecordNotFound.Error())
	require.Nil(t, getAccount)
}
