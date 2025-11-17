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

  if DB.Migrator().HasTable(&model.Event{}) {
    if err := DB.Exec(`ALTER TABLE events ADD COLUMN IF NOT EXISTS title varchar(255)`).Error; err != nil {
      return nil, err
    }
    if err := DB.Exec(`UPDATE events SET title = 'Untitled' WHERE title IS NULL`).Error; err != nil {
      return nil, err
    }
    if err := DB.Exec(`ALTER TABLE events ALTER COLUMN title SET NOT NULL`).Error; err != nil {
      return nil, err
    }
  }
  if err := DB.AutoMigrate(&model.User{}, &model.Event{}); err != nil {
    return nil, err
  }

	return DB, nil
}
