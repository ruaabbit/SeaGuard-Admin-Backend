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

// ListActivitiesForAdmin godoc
// @Summary 获取活动列表（管理员）
// @Description 获取所有志愿者活动的列表（需要管理员权限）
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.ActivitiesResponse "活动列表"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /admin/activities [get]
func (h *ActivityHandler) ListActivitiesForAdmin(c *gin.Context) {
    activities, err := h.service.GetAllActivities()
    if err != nil {
        c.JSON(500, models.ActivitiesResponse{
            Code:    500,
            Message: "获取活动列表失败",
            Data:    nil,
        })
        return
    }
    c.JSON(200, models.ActivitiesResponse{
        Code:    200,
        Message: "获取活动列表成功",
        Data:    activities,
    })
}

// ListAvailableActivities godoc
// @Summary 获取可报名活动列表（志愿者）
// @Description 获取所有可以报名的志愿者活动列表
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.ActivitiesResponse "活动列表"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities [get]
func (h *ActivityHandler) ListAvailableActivities(c *gin.Context) {
    activities, err := h.service.GetAvailableActivities()
    if err != nil {
        c.JSON(500, models.ActivitiesResponse{
            Code:    500,
            Message: "获取活动列表失败",
            Data:    nil,
        })
        return
    }
    c.JSON(200, models.ActivitiesResponse{
        Code:    200,
        Message: "获取可报名活动列表成功",
        Data:    activities,
    })
}

// CreateActivity godoc
// @Summary 创建新活动
// @Description 创建一个新的志愿者活动（需要管理员权限）
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param activity body models.Activity true "活动信息"
// @Success 201 {object} models.Activity "创建成功的活动信息"
// @Failure 400 {object} models.Response "请求参数无效"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities [post]
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	var activity models.Activity
if err := c.ShouldBindJSON(&activity); err != nil {
    c.JSON(400, models.Response{
        Error: "请求数据无效：" + err.Error(),
    })
    return
}

// 验证活动日期
if activity.Date.IsZero() {
    c.JSON(400, models.Response{
        Error: "活动日期不能为空",
    })
    return
}

if err := h.service.CreateActivity(&activity); err != nil {
    c.JSON(500, models.Response{
        Error: "创建活动失败：" + err.Error(),
    })
    return
}

c.JSON(201, models.Response{
    Message: "创建活动成功",
})
}

// UpdateActivity godoc
// @Summary 更新活动信息
// @Description 更新指定ID的活动信息（需要管理员权限）
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "活动ID"
// @Param activity body models.Activity true "活动信息"
// @Success 200 {object} models.Activity "更新后的活动信息"
// @Failure 400 {object} models.Response "无效的ID参数或请求数据"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities/{id} [put]
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	var activity models.Activity
if err := c.ShouldBindJSON(&activity); err != nil {
    c.JSON(400, models.Response{
        Error: "请求数据无效：" + err.Error(),
    })
    return
}

// 验证活动日期
if activity.Date.IsZero() {
    c.JSON(400, models.Response{
        Error: "活动日期不能为空",
    })
    return
}

id := uint(0)
if idParam := c.Param("id"); idParam != "" {
    if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
        id = uint(n)
    } else {
        c.JSON(400, models.Response{
            Error: "无效的ID参数",
        })
        return
    }
}

if err := h.service.UpdateActivity(id, &activity); err != nil {
    c.JSON(500, models.Response{
        Error: "更新活动失败：" + err.Error(),
    })
    return
}

c.JSON(200, models.Response{
    Message: "更新活动成功",
})
}

// DeleteActivity godoc
// @Summary 删除活动
// @Description 删除指定ID的活动（需要管理员权限）
// @Tags 活动管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "活动ID"
// @Success 200 {object} models.Response "活动删除成功"
// @Failure 400 {object} models.Response "无效的ID参数"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /activities/{id} [delete]
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
id := uint(0)
if idParam := c.Param("id"); idParam != "" {
    if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
        id = uint(n)
    } else {
        c.JSON(400, models.Response{
            Error: "无效的ID参数",
        })
        return
    }
}

if err := h.service.DeleteActivity(id); err != nil {
    c.JSON(500, models.Response{
        Error: "删除活动失败：" + err.Error(),
    })
    return
}

c.JSON(200, models.Response{
    Message: "活动删除成功",
})
}
