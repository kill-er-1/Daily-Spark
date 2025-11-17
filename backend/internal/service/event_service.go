package service

import (
    "context"
    "errors"
    "strings"
    "time"

    "github.com/cin/daily-spark/internal/model"
    "github.com/cin/daily-spark/internal/repository"
    "gorm.io/gorm"
)

type EventService struct {
	eventRepo *repository.EventRepository
	userRepo  *repository.UserRepository
}

func NewEventService(eventRepo *repository.EventRepository, userRepo *repository.UserRepository) *EventService {
	return &EventService{eventRepo: eventRepo, userRepo: userRepo}
}

func normalizeDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func (s *EventService) CreateEvent(ctx context.Context, userID string, content string, images []string, isPublic bool, eventDate *time.Time) (*model.Event, error) {
	userID = strings.TrimSpace(userID)
	content = strings.TrimSpace(content)

	if userID == "" {
		return nil, errors.New("user_id empty")
	}
	if content == "" {
		return nil, errors.New("content empty")
	}

	if _, err := s.userRepo.QueryUserByID(ctx, userID); err != nil {
		return nil, ErrUserNotFound
	}

	var date time.Time
	if eventDate == nil {
		date = normalizeDate(time.Now().In(time.Local))
	} else {
		date = normalizeDate(*eventDate)
	}

	e := &model.Event{
		UserID:    userID,
		Content:   content,
		Images:    images,
		IsPublic:  isPublic,
		EventDate: date,
	}

	created, err := s.eventRepo.CreateEvent(ctx, e)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, id string, title *string, content *string, images []string, isPublic *bool, eventDate *time.Time) (*model.Event, error) {
    id = strings.TrimSpace(id)
    if id == "" {
        return nil, errors.New("id empty")
    }

    e, err := s.eventRepo.QueryEventByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if title != nil {
        e.Title = strings.TrimSpace(*title)
    }
    if content != nil {
        e.Content = strings.TrimSpace(*content)
    }
    if isPublic != nil {
        e.IsPublic = *isPublic
    }
    if images != nil {
        e.Images = images
    }
    if eventDate != nil {
        e.EventDate = normalizeDate(*eventDate)
    }

    updated, err := s.eventRepo.UpdateEvent(ctx, e)
    if err != nil {
        return nil, err
    }
    return updated, nil
}

func (s *EventService) ListEventsByUserID(ctx context.Context, userID string) ([]*model.Event, error) {
    userID = strings.TrimSpace(userID)
    if userID == "" {
        return nil, errors.New("user_id empty")
    }
    if _, err := s.userRepo.QueryUserByID(ctx, userID); err != nil {
        return nil, ErrUserNotFound
    }
    return s.eventRepo.QueryEventsByUserID(ctx, userID)
}

func (s *EventService) DeleteEvent(ctx context.Context, id string) error {
    id = strings.TrimSpace(id)
    if id == "" {
        return errors.New("id empty")
    }
    return s.eventRepo.DeleteEventSoft(ctx, id)
}

var (
    ErrEventNotFound = errors.New("event not found")
)

func (s *EventService) AddEventTags(ctx context.Context, eventID string, names []string) (*model.Event, error) {
    eventID = strings.TrimSpace(eventID)
    if eventID == "" {
        return nil, errors.New("event_id empty")
    }
    if _, err := s.eventRepo.QueryEventByID(ctx, eventID); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrEventNotFound
        }
        return nil, err
    }
    tags, err := s.eventRepo.EnsureTagsByNames(ctx, names)
    if err != nil {
        return nil, err
    }
    if err := s.eventRepo.AddTagsToEvent(ctx, eventID, tags); err != nil {
        return nil, err
    }
    e, err := s.eventRepo.QueryEventByID(ctx, eventID)
    if err != nil {
        return nil, err
    }
    return e, nil
}

func (s *EventService) RemoveEventTags(ctx context.Context, eventID string, names []string) (*model.Event, error) {
    eventID = strings.TrimSpace(eventID)
    if eventID == "" {
        return nil, errors.New("event_id empty")
    }
    if _, err := s.eventRepo.QueryEventByID(ctx, eventID); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrEventNotFound
        }
        return nil, err
    }
    tags, err := s.eventRepo.EnsureTagsByNames(ctx, names)
    if err != nil {
        return nil, err
    }
    if err := s.eventRepo.RemoveTagsFromEvent(ctx, eventID, tags); err != nil {
        return nil, err
    }
    e, err := s.eventRepo.QueryEventByID(ctx, eventID)
    if err != nil {
        return nil, err
    }
    return e, nil
}

func (s *EventService) QueryEventsByTag(ctx context.Context, tagName string) ([]*model.Event, error) {
    tagName = strings.TrimSpace(tagName)
    if tagName == "" {
        return nil, errors.New("tag empty")
    }
    return s.eventRepo.QueryEventsByTagName(ctx, tagName)
}