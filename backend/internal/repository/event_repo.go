package repository

import (
	"context"
	"log"

	"github.com/cin/daily-spark/internal/model"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, e *model.Event) (*model.Event, error) {
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		log.Printf("[Repo.CreateEvent] db error: %v", err)
		return nil, err
	}
	return e, nil
}