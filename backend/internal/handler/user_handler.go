package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cin/daily-spark/internal/model"
	"github.com/cin/daily-spark/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// 文档用的请求体类型
type SignUpRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname,omitempty"`
	Password *string `json:"password,omitempty"`
}

type DeleteUserRequest struct {
	Account string `json:"account"`
}

// 用于文档与响应的用户视图类型（不含敏感与 ORM 私有字段）
type UserVO struct {
	ID        string  `json:"id"`
	Account   string  `json:"account"`
	Nickname  *string `json:"nickname,omitempty"`
	IsAdmin   bool    `json:"is_admin"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func toVO(u *model.User) *UserVO {
	if u == nil {
		return nil
	}
	return &UserVO{
		ID:        u.ID,
		Account:   u.Account,
		Nickname:  u.Nickname,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// SignUp godoc
// @Summary 用户注册
// @Description 使用账号密码创建用户（返回创建的用户）
// @Tags users
// @Accept json
// @Produce json
// @Param body body SignUpRequest true "注册参数"
// @Success 200 {object} UserVO
// @Failure 400 {string} string "bad request"
// @Failure 409 {string} string "account exists"
// @Failure 500 {string} string "internal error"
// @Router /users/signup [post]
func (s *UserHandler) SignUp(c *gin.Context) {
	var req struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SignUp] bind json error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SignUp] try create user, account=%s", strings.TrimSpace(req.Account))
	user, err := s.service.UserSignUp(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrAccountExists) {
			log.Printf("[SignUp] account exists: account=%s", strings.TrimSpace(req.Account))
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Printf("[SignUp] internal error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("[SignUp] user created: id=%s account=%s", user.ID, user.Account)
	c.JSON(http.StatusOK, gin.H{"message": "user created", "user": user})
}

// SignIn godoc
// @Summary 用户登录
// @Description 使用账号密码登录（返回用户信息）
// @Tags users
// @Accept json
// @Produce json
// @Param body body SignInRequest true "登录参数"
// @Success 200 {object} UserVO
// @Failure 400 {string} string "bad request"
// @Failure 401 {string} string "password not match"
// @Failure 404 {string} string "account not exists"
// @Failure 500 {string} string "internal error"
// @Router /users/signin [post]
func (s *UserHandler) SignIn(c *gin.Context) {
	var req struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SignIn] bind json error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SignIn] try sign in, account=%s", strings.TrimSpace(req.Account))
	user, err := s.service.UserSignIn(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotExists) {
			log.Printf("[SignIn] account not exists: account=%s", strings.TrimSpace(req.Account))
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrPasswordNotMatch) {
			log.Printf("[SignIn] password not match: account=%s", strings.TrimSpace(req.Account))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			log.Printf("[SignIn] internal error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("[SignIn] user signed in: id=%s account=%s", user.ID, user.Account)
	c.JSON(http.StatusOK, gin.H{"message": "user signed in", "user": user})
}

// QueryAllUsers godoc
// @Summary 查询所有用户
// @Tags users
// @Produce json
// @Success 200 {array} UserVO
// @Failure 500 {string} string "internal error"
// @Router /users/query [get]
func (s *UserHandler) QueryAllUsers(c *gin.Context) {
	users, err := s.service.QueryAllUsers(c.Request.Context())
	if err != nil {
		log.Printf("[QueryAllUsers] error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[QueryAllUsers] success: count=%d", len(users))
	c.JSON(http.StatusOK, gin.H{"message": "users queried", "users": users})
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 根据 id 更新昵称或密码
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param body body UpdateUserRequest true "更新参数"
// @Success 200 {object} UserVO
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "user not found"
// @Failure 500 {string} string "internal error"
// @Router /users/update/{id} [post]
func (s *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Printf("[UpdateUser] empty id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}

	var req struct {
		Nickname *string `json:"nickname,omitempty"`
		Password *string `json:"password,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateUser] bind json error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[UpdateUser] try update user: id=%s nickname_set=%t password_set=%t", id, req.Nickname != nil, req.Password != nil)
	user, err := s.service.UpdateUser(c.Request.Context(), id, req.Nickname, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Printf("[UpdateUser] user not found: id=%s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("[UpdateUser] internal error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("[UpdateUser] user updated: id=%s", user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "user updated", "user": user})
}
	// DeleteUser godoc
	// @Summary 删除用户（管理员）
	// @Description 管理员通过 id 软删除用户
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param id path string true "用户ID"
	// @Param body body DeleteUserRequest true "操作者账号"
	// @Success 200 {string} string "user deleted"
	// @Failure 400 {string} string "bad request"
	// @Failure 403 {string} string "permission denied"
	// @Failure 404 {string} string "user not found"
	// @Failure 500 {string} string "internal error"
	// @Router /users/delete/{id} [post]
func (s *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Printf("[DeleteUser] empty id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}

	var req DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[DeleteUser] bind json error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account := strings.TrimSpace(req.Account)
	if account == "" {
		log.Printf("[DeleteUser] empty operator account")
		c.JSON(http.StatusBadRequest, gin.H{"error": "operator account is empty"})
		return
	}

	log.Printf("[DeleteUser] operator=%s try delete id=%s", account, id)
	if err := s.service.DeleteUser(c.Request.Context(), id, account); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Printf("[DeleteUser] user not found: id=%s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrPermissionDenied) {
			log.Printf("[DeleteUser] permission denied: operator=%s id=%s", account, id)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			log.Printf("[DeleteUser] internal error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("[DeleteUser] user deleted: id=%s by operator=%s", id, account)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}