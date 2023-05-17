package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbName)

	ctx := context.Background()
	timeout := 5 * time.Second

	// 使用 context.WithTimeout 創建一個帶有超時的子 context
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)

	var gormDB *gorm.DB
	var err error

	go func() {
		for {
			select {
			case <-ctx.Done():
				done <- fmt.Errorf("db connection time out")
				break
			default:
				gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
					Logger:         logger.Default.LogMode(logger.Info),
					TranslateError: true,
				})
				if err == nil {
					done <- nil
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, err
		}
	}

	return gormDB, nil
}
