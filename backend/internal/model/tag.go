package model

import (
    "gorm.io/gorm"
)

type Tag struct {
    ID   string `gorm:"type:uuid;default:gen_random_uuid();primaryKey;->" json:"id"`
    Name string `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
}

func (Tag) TableName() string { return "tags" }

type EventTag struct {
    EventID string `gorm:"type:uuid;primaryKey;not null"`
    TagID   string `gorm:"type:uuid;primaryKey;not null"`
    Event   Event  `gorm:"foreignKey:EventID;references:ID;constraint:OnDelete:CASCADE"`
    Tag     Tag    `gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (EventTag) TableName() string { return "event_tags" }