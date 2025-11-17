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
    Title     string   `json:"title"`
    Content   string   `json:"content"`
    Images    []string `json:"images,omitempty"`
    IsPublic  bool     `json:"is_public"`
    EventDate string   `json:"event_date"`
    CreatedAt string   `json:"created_at"`
    UpdatedAt string   `json:"updated_at"`
}

// swagger response models
type EventUpdateResponse struct {
    Message string  `json:"message"`
    Event   EventVO `json:"event"`
}

type EventListResponse struct {
    Events []EventVO `json:"events"`
}

type SimpleMessageResponse struct {
    Message string `json:"message"`
}

func eventToVO(e *model.Event) *EventVO {
    if e == nil {
        return nil
    }
    return &EventVO{
        ID:        e.ID,
        UserID:    e.UserID,
        Title:     e.Title,
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

// UpdateEvent
type UpdateEventRequest struct {
    Title     *string  `json:"title,omitempty"`
    Content   *string  `json:"content,omitempty"`
    Images    []string `json:"images,omitempty"`
    IsPublic  *bool    `json:"is_public,omitempty"`
    EventDate *string  `json:"event_date,omitempty"` // YYYY-MM-DD
}

// UpdateEvent godoc
// @Summary 编辑事件
// @Description 根据事件ID更新事件，支持部分字段更新
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "事件ID"
// @Param body body UpdateEventRequest true "更新内容"
// @Success 200 {object} EventUpdateResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /events/update/{id} [post]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
    id := strings.TrimSpace(c.Param("id"))
    var req UpdateEventRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var dt *time.Time
    if req.EventDate != nil && *req.EventDate != "" {
        parsed, err := time.ParseInLocation("2006-01-02", *req.EventDate, time.Local)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_date, expected YYYY-MM-DD"})
            return
        }
        dt = &parsed
    }

    e, err := h.service.UpdateEvent(c.Request.Context(), id, req.Title, req.Content, req.Images, req.IsPublic, dt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "event updated", "event": eventToVO(e)})
}

// Query by user_id
// QueryEventsByUserID godoc
// @Summary 查询用户下的事件
// @Description 根据用户ID返回事件列表（排除已软删除）
// @Tags events
// @Accept json
// @Produce json
// @Param user_id query string true "用户ID"
// @Success 200 {object} EventListResponse
// @Failure 404 {string} string "user not found"
// @Failure 500 {string} string "internal error"
// @Router /events/query [get]
func (h *EventHandler) QueryEventsByUserID(c *gin.Context) {
    userID := strings.TrimSpace(c.Query("user_id"))
    list, err := h.service.ListEventsByUserID(c.Request.Context(), userID)
    if err != nil {
        if errors.Is(err, service.ErrUserNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    vos := make([]*EventVO, 0, len(list))
    for _, e := range list {
        vos = append(vos, eventToVO(e))
    }
    c.JSON(http.StatusOK, gin.H{"events": vos})
}

// DeleteEvent
type DeleteEventRequest struct {}

// DeleteEvent godoc
// @Summary 删除事件（软删除）
// @Description 软删除事件记录
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "事件ID"
// @Param body body DeleteEventRequest false "空 JSON 包装（可选）"
// @Success 200 {object} SimpleMessageResponse
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal error"
// @Router /events/delete/{id} [post]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
    id := strings.TrimSpace(c.Param("id"))
    var req DeleteEventRequest
    if c.Request.ContentLength > 0 {
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    }
    if err := h.service.DeleteEvent(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}
