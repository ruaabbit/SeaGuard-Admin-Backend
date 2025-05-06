package handlers

import (
"github.com/gin-gonic/gin"
"net/http"
"seaguard-admin-backend/models"
"seaguard-admin-backend/service"
"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary 用户注册
// @Description 创建新用户账号。如果role为volunteer，则同时创建志愿者信息
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册信息。当role为volunteer时，需要提供name、phone、email、address等志愿者信息"
// @Success 200 {object} models.Response "注册成功"
// @Failure 400 {object} models.Response "请求参数无效: 1. 用户名已存在 2. 必填字段缺失 3. role为volunteer时未提供志愿者信息"
// @Example {
//   "request": {
//     "username": "zhangsan",
//     "password": "123456",
//     "role": "volunteer",
//     "name": "张三",
//     "phone": "13800138000",
//     "email": "zhangsan@example.com",
//     "address": "北京市海淀区"
//   }
// }
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if req.Role != "admin" && req.Role != "volunteer" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户角色"})
        return
    }

    err := h.userService.Register(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// @Summary 用户登录
// @Description 用户登录并获取认证token
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录信息"
// @Success 200 {object} models.LoginResponse "登录成功"
// @Failure 400 {object} models.Response "请求参数无效"
// @Failure 401 {object} models.Response "用户名或密码错误"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	user, token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// @Summary 修改密码
// @Description 修改用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.ChangePasswordRequest true "密码修改信息"
// @Success 200 {object} models.Response "密码修改成功"
// @Failure 400 {object} models.Response "请求参数无效或密码错误"
// @Router /auth/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	userID := c.GetUint("userID")
	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// @Summary 获取用户列表
// @Description 获取所有用户信息（仅管理员可用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.UsersResponse "用户列表"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// @Summary 更新用户状态
// @Description 更新指定用户的状态（仅管理员可用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body models.StatusUpdateRequest true "状态信息"
// @Success 200 {object} models.Response "状态更新成功"
// @Failure 400 {object} models.Response "无效的用户ID或状态"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	if err := h.userService.UpdateStatus(uint(userID), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户状态更新成功"})
}

// @Summary 删除用户
// @Description 删除指定用户（仅管理员可用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} models.Response "用户删除成功"
// @Failure 400 {object} models.Response "无效的用户ID"
// @Failure 403 {object} models.Response "无权限访问"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}
