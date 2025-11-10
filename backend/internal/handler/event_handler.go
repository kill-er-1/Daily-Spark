package handler

import (
	"errors"
	"log"
	"net/http"
	"time"
	// 新增导入
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cin/daily-spark/internal/model"
	"github.com/cin/daily-spark/internal/service"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(service *service.EventService) *EventHandler {
	return &EventHandler{service: service}
}

// 请求体
type CreateEventRequest struct {
	UserID    string   `json:"user_id"`
	Content   string   `json:"content"`
	Images    []string `json:"images,omitempty"`
	IsPublic  bool     `json:"is_public"`
	EventDate *string  `json:"event_date,omitempty"` // YYYY-MM-DD，可选
}

// 响应 VO（避免 DeletedAt）
type EventVO struct {
	ID        string   `json:"id"`
	UserID    string   `json:"user_id"`
	Content   string   `json:"content"`
	Images    []string `json:"images,omitempty"`
	IsPublic  bool     `json:"is_public"`
	EventDate string   `json:"event_date"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func eventToVO(e *model.Event) *EventVO {
	if e == nil {
		return nil
	}
	return &EventVO{
		ID:        e.ID,
		UserID:    e.UserID,
		Content:   e.Content,
		Images:    e.Images,
		IsPublic:  e.IsPublic,
		EventDate: e.EventDate.Format("2006-01-02"),
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.Format(time.RFC3339),
	}
}

// CreateEvent godoc
// @Summary 创建快乐事件
// @Description 为指定日期（默认当天）创建一条快乐事件；图片上传为 multipart/form-data（单文件）
// @Tags events
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData string true "用户ID"
// @Param content formData string true "纯文本内容"
// @Param is_public formData boolean false "是否公开（true/false）"
// @Param event_date formData string false "事件日期(YYYY-MM-DD)"
// @Param image formData file false "事件图片（单文件）"
// @Success 200 {object} EventVO
// @Failure 400 {string} string "bad request"
// @Failure 404 {string} string "user not found"
// @Failure 500 {string} string "internal error"
// @Router /events/create [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {

	// 删除 JSON 绑定，改为从 multipart/form-data 表单读取
	userID := strings.TrimSpace(c.PostForm("user_id"))
	content := strings.TrimSpace(c.PostForm("content"))
	isPublicStr := strings.TrimSpace(c.DefaultPostForm("is_public", "false"))
	isPublic := strings.EqualFold(isPublicStr, "true")
	eventDateStr := strings.TrimSpace(c.PostForm("event_date"))

	if userID == "" || content == "" {
		log.Printf("[CreateEvent] empty user_id or content")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or content empty"})
		return
	}

	// 解析日期（可选，默认当天）
	var dt *time.Time
	if eventDateStr != "" {
		parsed, err := time.ParseInLocation("2006-01-02", eventDateStr, time.Local)
		if err != nil {
			log.Printf("[CreateEvent] bad event_date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_date, expected YYYY-MM-DD"})
			return
		}
		dt = &parsed
	}

	// 处理单文件上传：字段名 image
	uploadDir := filepath.Join("static", "uploads", "events", userID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("[CreateEvent] mkdir error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var imagePaths []string
	if f, err := c.FormFile("image"); err == nil && f != nil {
		ext := filepath.Ext(f.Filename)
		base := strings.TrimSuffix(f.Filename, ext)
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), base, ext)
		fullPath := filepath.Join(uploadDir, filename)
		if err := c.SaveUploadedFile(f, fullPath); err != nil {
			log.Printf("[CreateEvent] save file error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		publicPath := filepath.ToSlash(filepath.Join("/static", "uploads", "events", userID, filename))
		imagePaths = append(imagePaths, publicPath)
	}

	// 调用服务层
	e, err := h.service.CreateEvent(c.Request.Context(), userID, content, imagePaths, isPublic, dt)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			log.Printf("[CreateEvent] user not found: user_id=%s", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("[CreateEvent] internal error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	log.Printf("[CreateEvent] created: id=%s user_id=%s images=%d", e.ID, e.UserID, len(imagePaths))
	c.JSON(http.StatusOK, eventToVO(e))
}
