package handlers

import (
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/service"
	"strconv"

	"github.com/gin-gonic/gin"
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

// Register godoc
// @Summary 报名活动
// @Description 志愿者报名参加活动
// @Tags 报名管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "活动ID"
// @Param registration body models.RegistrationRequest true "报名信息"
// @Success 201 {object} models.Response "报名成功"
// @Failure 400 {object} models.Response "无效的活动ID或报名信息"
// @Failure 401 {object} models.Response "未登录"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 409 {object} models.Response "已经报名过该活动"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities/{id}/register [post]
func (h *RegistrationHandler) Register(c *gin.Context) {
	// 获取活动ID
	activityID := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			activityID = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的活动ID参数"})
			return
		}
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "未登录"})
		return
	}

	// 解析请求体
	var req models.RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 创建报名记录
	registration := &models.Registration{
		ActivityID:       activityID,
		Name:             req.Name,
		Phone:            req.Phone,
		IDCard:           req.IDCard,
		Email:            req.Email,
		EmergencyContact: req.EmergencyContact,
		EmergencyPhone:   req.EmergencyPhone,
		Status:           "pending",
	}

	if err := h.service.CreateRegistration(userID.(uint), registration); err != nil {
		if err.Error() == "already registered" {
			c.JSON(409, gin.H{"error": "已经报名过该活动"})
			return
		}
		c.JSON(500, gin.H{"error": "报名失败"})
		return
	}

	c.JSON(201, gin.H{"message": "报名成功"})
}

// GetMyRegistration godoc
// @Summary 查询个人报名状态
// @Description 查询当前用户在指定活动的报名状态
// @Tags 报名管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "活动ID"
// @Success 200 {object} models.Registration "报名信息"
// @Failure 400 {object} models.Response "无效的活动ID"
// @Failure 401 {object} models.Response "未登录"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 404 {object} models.Response "未找到报名记录"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities/{id}/registration [get]
func (h *RegistrationHandler) GetMyRegistration(c *gin.Context) {
	// 获取活动ID
	activityID := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			activityID = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的活动ID参数"})
			return
		}
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "未登录"})
		return
	}

	// 查询报名记录
	registration, err := h.service.GetUserRegistration(userID.(uint), activityID)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "未找到报名记录"})
			return
		}
		c.JSON(500, gin.H{"error": "获取报名记录失败"})
		return
	}

	c.JSON(200, registration)
}
