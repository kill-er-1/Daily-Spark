package database

import (
	"github.com/cin/daily-spark/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=postgres user=postgres password=postgres dbname=daily_spark port=5432 sslmode=disable TimeZone=Asia/Shanghai"

func InitDB() (*gorm.DB, error) {
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//postgres 生成uuid插件
	if err := DB.Exec(`create extension if not exists pgcrypto`).Error; err != nil {
		return nil, err
	}

  //数据库迁移
  if err := DB.AutoMigrate(&model.User{}); err != nil {
    return nil, err
  }

	return DB, nil
}
