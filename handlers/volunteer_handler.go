package handlers

import (
	"github.com/gin-gonic/gin"
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/service"
	"strconv"
)

// VolunteerHandler 志愿者处理器结构
type VolunteerHandler struct {
	service service.VolunteerService
}

// NewVolunteerHandler 创建志愿者处理器实例
func NewVolunteerHandler(service service.VolunteerService) *VolunteerHandler {
	return &VolunteerHandler{
		service: service,
	}
}

// ListVolunteers godoc
// @Summary 获取志愿者列表
// @Description 获取所有志愿者的列表（需要管理员权限）
// @Tags 志愿者管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.VolunteersResponse "志愿者列表"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /volunteers [get]
func (h *VolunteerHandler) ListVolunteers(c *gin.Context) {
	volunteers, err := h.service.GetAllVolunteers()
	if err != nil {
		c.JSON(500, gin.H{"error": "获取志愿者列表失败"})
		return
	}
	c.JSON(200, volunteers)
}

// CreateVolunteer godoc
// @Summary 创建新志愿者
// @Description 创建一个新的志愿者信息（需要管理员权限）
// @Tags 志愿者管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param volunteer body models.Volunteer true "志愿者信息"
// @Success 201 {object} models.Volunteer "创建成功的志愿者信息"
// @Failure 400 {object} models.Response "请求参数无效"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /volunteers [post]
func (h *VolunteerHandler) CreateVolunteer(c *gin.Context) {
	var volunteer models.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateVolunteer(&volunteer); err != nil {
		c.JSON(500, gin.H{"error": "创建志愿者失败"})
		return
	}

	c.JSON(201, volunteer)
}

// UpdateVolunteer godoc
// @Summary 更新志愿者信息
// @Description 更新指定ID的志愿者信息（需要管理员权限）
// @Tags 志愿者管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "志愿者ID"
// @Param volunteer body models.Volunteer true "志愿者信息"
// @Success 200 {object} models.Volunteer "更新后的志愿者信息"
// @Failure 400 {object} models.Response "无效的ID参数或请求数据"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /volunteers/{id} [put]
func (h *VolunteerHandler) UpdateVolunteer(c *gin.Context) {
	var volunteer models.Volunteer
	if err := c.ShouldBindJSON(&volunteer); err != nil {
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

	if err := h.service.UpdateVolunteer(id, &volunteer); err != nil {
		c.JSON(500, gin.H{"error": "更新志愿者失败"})
		return
	}

	c.JSON(200, volunteer)
}

// DeleteVolunteer godoc
// @Summary 删除志愿者
// @Description 删除指定ID的志愿者（需要管理员权限）
// @Tags 志愿者管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "志愿者ID"
// @Success 200 {object} models.Response "志愿者删除成功"
// @Failure 400 {object} models.Response "无效的ID参数"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /volunteers/{id} [delete]
func (h *VolunteerHandler) DeleteVolunteer(c *gin.Context) {
	id := uint(0)
	if idParam := c.Param("id"); idParam != "" {
		if n, err := strconv.ParseUint(idParam, 10, 32); err == nil {
			id = uint(n)
		} else {
			c.JSON(400, gin.H{"error": "无效的ID参数"})
			return
		}
	}

	if err := h.service.DeleteVolunteer(id); err != nil {
		c.JSON(500, gin.H{"error": "删除志愿者失败"})
		return
	}
c.JSON(200, gin.H{"message": "志愿者删除成功"})
}

// UpdateMyInfo godoc
// @Summary 更新个人志愿者信息
// @Description 已登录的志愿者用户更新自己的个人信息
// @Tags 志愿者
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param info body models.UpdateVolunteerInfoRequest true "志愿者个人信息 (姓名、电话、邮箱、地址)"
// @Success 200 {object} models.Response "更新成功"
// @Failure 400 {object} models.Response "请求参数无效：1. 必填字段缺失 2. 邮箱格式错误"
// @Failure 401 {object} models.Response "未登录"
// @Failure 403 {object} models.Response "无权限访问：非志愿者用户"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Example {
//   "request": {
//     "name": "张三",
//     "phone": "13800138000",
//     "email": "zhangsan@example.com",
//     "address": "北京市海淀区"
//   }
// }
// @Router /volunteer/my-info [put]
func (h *VolunteerHandler) UpdateMyInfo(c *gin.Context) {
    // 从上下文获取当前用户ID
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(401, gin.H{"error": "未登录"})
        return
    }

    // 验证请求体
    var req models.UpdateVolunteerInfoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 更新志愿者信息
    if err := h.service.UpdateVolunteerInfo(uint(userID.(float64)), &req); err != nil {
        c.JSON(500, gin.H{"error": "更新个人信息失败"})
        return
    }

    c.JSON(200, gin.H{"message": "个人信息更新成功"})
}
