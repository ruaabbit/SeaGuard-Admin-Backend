package models

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	// User账号信息
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"your_password"`
	Role     string `json:"role" binding:"required" example:"volunteer" enums:"admin,volunteer"`

	// Volunteer个人信息（当role为volunteer时必填）
	Name    string `json:"name" binding:"required_if=Role volunteer" example:"张三"`
	Phone   string `json:"phone" binding:"required_if=Role volunteer" example:"13800138000"`
	Email   string `json:"email" binding:"required_if=Role volunteer,email" example:"zhangsan@example.com"`
	Address string `json:"address" binding:"required_if=Role volunteer" example:"北京市海淀区"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"your_password"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"old_password"`
	NewPassword string `json:"new_password" binding:"required" example:"new_password"`
}

// StatusUpdateRequest 状态更新请求已在response.go中定义

// UpdateVolunteerInfoRequest 更新志愿者信息请求
type UpdateVolunteerInfoRequest struct {
	Name    string `json:"name" binding:"required" example:"张三"`
	Phone   string `json:"phone" binding:"required" example:"13800138000"`
	Email   string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Address string `json:"address" binding:"required" example:"北京市海淀区"`
}
