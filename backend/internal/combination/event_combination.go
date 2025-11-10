package combination

import (
	"errors"

	"github.com/cin/daily-spark/internal/handler"
	"github.com/cin/daily-spark/internal/repository"
	"github.com/cin/daily-spark/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventModule struct {
	DB      *gorm.DB
	Repo    *repository.EventRepository
	UserRepo *repository.UserRepository
	Service *service.EventService
	Handler *handler.EventHandler
}

func NewEventModule(db *gorm.DB) *EventModule {
	eventRepo := repository.NewEventRepository(db)
	userRepo := repository.NewUserRepository(db)
	svc := service.NewEventService(eventRepo, userRepo)
	h := handler.NewEventHandler(svc)

	return &EventModule{
		DB:       db,
		Repo:     eventRepo,
		UserRepo: userRepo,
		Service:  svc,
		Handler:  h,
	}
}

func BuildEventModule(db *gorm.DB) (*EventModule, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return NewEventModule(db), nil
}

func RegisterEventRoutes(r *gin.Engine, h *handler.EventHandler) {
	v1 := r.Group("/api/v1/events")
	v1.POST("/create", h.CreateEvent)
}