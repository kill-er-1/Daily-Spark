package model

import (
    "time"

    "gorm.io/gorm"
)

type User struct {
    ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey;->" json:"id"`
    Account   string         `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_account;<-:create" json:"account"`
    Password  string         `gorm:"type:varchar(255);not null;<-:create,update" json:"-"`
    Nickname  *string        `gorm:"type:varchar(255);<-" json:"nickname,omitempty"`
    IsAdmin   bool           `gorm:"not null;default:false;->;<-:create" json:"is_admin"`
    CreatedAt time.Time      `gorm:"not null;default:now();->" json:"created_at"`
    UpdatedAt time.Time      `gorm:"not null;default:now();->" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
    return "users"
}