package model

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
    ID         string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey;->" json:"id"`
    UserID     string         `gorm:"type:uuid;not null;index:idx_events_user_id" json:"user_id"`
    Title      string         `gorm:"type:varchar(255);not null" json:"title"`
    Content    string         `gorm:"type:text;not null" json:"content"`
    Images     []string       `gorm:"type:jsonb;serializer:json" json:"images,omitempty"`
    EventDate  time.Time      `gorm:"type:date;not null;column:event_date;index:idx_events_event_date" json:"event_date"`
    IsPublic   bool           `gorm:"not null;default:false;<-:create,update" json:"is_public"`
    IsFeatured bool           `gorm:"not null;default:false;<-:create,update" json:"is_featured"`
    Tags       []Tag          `gorm:"many2many:event_tags" json:"tags,omitempty"`
    CreatedAt  time.Time      `gorm:"not null;default:now();->" json:"created_at"`
    UpdatedAt  time.Time      `gorm:"not null;default:now();->" json:"updated_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index:idx_events_deleted_at" json:"deleted_at,omitempty"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Event) TableName() string {
    return "events"
}