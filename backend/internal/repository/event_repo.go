package repository

import (
    "context"
    "log"
    "strings"

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
    if err := r.db.WithContext(ctx).
        Where("id = ? and deleted_at is null", id).
        Preload("Tags").
        First(&e).Error; err != nil {
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

func (r *EventRepository) EnsureTagsByNames(ctx context.Context, names []string) ([]model.Tag, error) {
    out := make([]model.Tag, 0, len(names))
    for _, n := range names {
        name := strings.TrimSpace(n)
        if name == "" { continue }
        var t model.Tag
        if err := r.db.WithContext(ctx).Where("name = ?", name).First(&t).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                t = model.Tag{Name: name}
                if err := r.db.WithContext(ctx).Create(&t).Error; err != nil {
                    log.Printf("[Repo.EnsureTagsByNames] create tag error: %v", err)
                    return nil, err
                }
            } else {
                return nil, err
            }
        }
        out = append(out, t)
    }
    return out, nil
}

func (r *EventRepository) AddTagsToEvent(ctx context.Context, eventID string, tags []model.Tag) error {
    e := &model.Event{ID: eventID}
    if err := r.db.WithContext(ctx).Model(e).Association("Tags").Append(&tags); err != nil {
        log.Printf("[Repo.AddTagsToEvent] assoc error: %v", err)
        return err
    }
    return nil
}

func (r *EventRepository) RemoveTagsFromEvent(ctx context.Context, eventID string, tags []model.Tag) error {
    e := &model.Event{ID: eventID}
    if err := r.db.WithContext(ctx).Model(e).Association("Tags").Delete(&tags); err != nil {
        log.Printf("[Repo.RemoveTagsFromEvent] assoc error: %v", err)
        return err
    }
    return nil
}

func (r *EventRepository) QueryEventsByTagName(ctx context.Context, tagName string) ([]*model.Event, error) {
    var list []*model.Event
    err := r.db.WithContext(ctx).
        Joins("JOIN event_tags ON event_tags.event_id = events.id").
        Joins("JOIN tags ON tags.id = event_tags.tag_id").
        Where("tags.name = ? AND events.deleted_at IS NULL", tagName).
        Order("events.event_date desc, events.created_at desc").
        Preload("Tags").
        Find(&list).Error
    if err != nil {
        log.Printf("[Repo.QueryEventsByTagName] db error: %v", err)
        return nil, err
    }
    return list, nil
}