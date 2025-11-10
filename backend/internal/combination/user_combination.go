package combination

import (
	"errors"

	"github.com/cin/daily-spark/internal/handler"
	"github.com/cin/daily-spark/internal/repository"
	"github.com/cin/daily-spark/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserModule struct {
	DB      *gorm.DB
	Repo    *repository.UserRepository
	Service *service.UserService
	Handler *handler.UserHandler
}

func NewUserModule(db *gorm.DB) *UserModule {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	return &UserModule{
		DB:      db,
		Repo:    repo,
		Service: svc,
		Handler: h,
	}
}

func BuildUserModule(db *gorm.DB) (*UserModule, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return NewUserModule(db), nil
}

func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler) {
	v1 := r.Group("/api/v1/users")
	v1.POST("/signup", h.SignUp)
	v1.POST("/signin", h.SignIn)
	v1.GET("/query", h.QueryAllUsers)
	v1.POST("/update/:id", h.UpdateUser)
	v1.POST("/delete/:id", h.DeleteUser)
}