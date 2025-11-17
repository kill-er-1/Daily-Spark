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

func (r *EventRepository) QueryEventByID(ctx context.Context,id string )(*model.Event,error) {
	var e model.Event
	if err := r.db.WithContext(ctx).Where("id = ? and deleted_at is null", id).First(&e).Error; err != nil {
		log.Printf("[Repo.QueryEventByID] db error: %v", err)
		return nil, err
	}
	return &e, nil	
}

func (r *EventRepository) UpdateEvent(ctx context.Context, e *model.Event) (*model.Event, error) {
    if err := r.db.WithContext(ctx).Save(e).Error; err != nil {
        log.Printf("[Repo.UpdateEvent] db error: %v", err)
        return nil, err
    }
    return e, nil
}

func (r *EventRepository) QueryEventsByUserID(ctx context.Context, userID string) ([]*model.Event, error) {
    var list []*model.Event
    if err := r.db.WithContext(ctx).
        Where("user_id = ? and deleted_at is null", userID).
        Order("event_date desc, created_at desc").
        Find(&list).Error; err != nil {
        log.Printf("[Repo.QueryEventsByUserID] db error: %v", err)
        return nil, err
    }
    return list, nil
}

func (r *EventRepository) DeleteEventSoft(ctx context.Context, id string) error {
    if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Event{}).Error; err != nil {
        log.Printf("[Repo.DeleteEventSoft] db error: %v", err)
        return err
    }
    return nil
}