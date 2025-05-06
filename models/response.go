package models

// Response 基础响应结构
type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    User   `json:"user"`
}

// UsersResponse 用户列表响应结构
type UsersResponse struct {
	Users []User `json:"users"`
}

// ActivitiesResponse 活动列表响应结构
type ActivitiesResponse struct {
	Activities []Activity `json:"activities"`
}

// VolunteersResponse 志愿者列表响应结构
type VolunteersResponse struct {
	Volunteers []Volunteer `json:"volunteers"`
}

// RegistrationsResponse 报名记录列表响应结构
type RegistrationsResponse struct {
	Registrations []Registration `json:"registrations"`
}

// StatusUpdateRequest 状态更新请求结构
type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}
