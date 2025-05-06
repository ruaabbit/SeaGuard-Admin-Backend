package handlers

import (
	"github.com/gin-gonic/gin"
	"seaguard-admin-backend/service"
	"strconv"
)

// RegistrationHandler 报名记录处理器结构
type RegistrationHandler struct {
	service service.RegistrationService
}

// NewRegistrationHandler 创建报名记录处理器实例
func NewRegistrationHandler(service service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{
		service: service,
	}
}

// ListActivityRegistrations godoc
// @Summary 获取活动报名列表
// @Description 获取指定活动的所有报名记录（需要志愿者权限）
// @Tags 报名管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "活动ID"
// @Success 200 {object} models.RegistrationsResponse "报名记录列表"
// @Failure 400 {object} models.Response "无效的活动ID"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities/{id}/registrations [get]
func (h *RegistrationHandler) ListActivityRegistrations(c *gin.Context) {
	id := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			id = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的活动ID参数"})
			return
		}
	}

	registrations, err := h.service.GetActivityRegistrations(id)
	if err != nil {
		c.JSON(500, gin.H{"error": "获取报名列表失败"})
		return
	}

	c.JSON(200, registrations)
}

// UpdateRegistrationStatus godoc
// @Summary 更新报名状态
// @Description 更新指定报名记录的状态（需要志愿者权限）
// @Tags 报名管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "报名ID"
// @Param status body models.StatusUpdateRequest true "状态信息（可选值：pending待审核、approved已通过、rejected已拒绝）"
// @Success 200 {object} models.Response "状态更新成功"
// @Failure 400 {object} models.Response "无效的报名ID或状态值"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /registrations/{id}/status [put]
func (h *RegistrationHandler) UpdateRegistrationStatus(c *gin.Context) {
	id := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			id = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的报名ID参数"})
			return
		}
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateRegistrationStatus(id, statusUpdate.Status); err != nil {
		c.JSON(500, gin.H{"error": "更新报名状态失败"})
		return
	}

	c.JSON(200, gin.H{"message": "报名状态更新成功"})
}
