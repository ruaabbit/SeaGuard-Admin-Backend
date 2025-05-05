package handlers

import (
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ActivityHandler 活动处理器结构
type ActivityHandler struct {
	service service.ActivityService
}

// NewActivityHandler 创建活动处理器实例
func NewActivityHandler(service service.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

// ListActivities godoc
// @Summary 获取活动列表
// @Description 获取所有志愿者活动的列表
// @Tags 活动管理
// @Accept json
// @Produce json
// @Success 200 {array} models.Activity
// @Failure 500 {object} map[string]string
// @Router /activities [get]
func (h *ActivityHandler) ListActivities(c *gin.Context) {
	activities, err := h.service.GetAllActivities()
	if err != nil {
		c.JSON(500, gin.H{"error": "获取活动列表失败"})
		return
	}
	c.JSON(200, activities)
}

// CreateActivity godoc
// @Summary 创建新活动
// @Description 创建一个新的志愿者活动
// @Tags 活动管理
// @Accept json
// @Produce json
// @Param activity body models.Activity true "活动信息"
// @Success 201 {object} models.Activity
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /activities [post]
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateActivity(&activity); err != nil {
		c.JSON(500, gin.H{"error": "创建活动失败"})
		return
	}

	c.JSON(201, activity)
}

// UpdateActivity godoc
// @Summary 更新活动信息
// @Description 更新指定ID的活动信息
// @Tags 活动管理
// @Accept json
// @Produce json
// @Param id path int true "活动ID"
// @Param activity body models.Activity true "活动信息"
// @Success 200 {object} models.Activity
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /activities/{id} [put]
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	var activity models.Activity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			id = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的ID参数"})
			return
		}
	}
	if err := h.service.UpdateActivity(id, &activity); err != nil {
		c.JSON(500, gin.H{"error": "更新活动失败"})
		return
	}

	c.JSON(200, activity)
}

// DeleteActivity godoc
// @Summary 删除活动
// @Description 删除指定ID的活动
// @Tags 活动管理
// @Accept json
// @Produce json
// @Param id path int true "活动ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /activities/{id} [delete]
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	id := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			id = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的ID参数"})
			return
		}
	}
	if err := h.service.DeleteActivity(id); err != nil {
		c.JSON(500, gin.H{"error": "删除活动失败"})
		return
	}
	c.JSON(200, gin.H{"message": "活动删除成功"})
}
