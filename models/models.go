package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	Role      string    `json:"role"` // admin或volunteer
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Activity 活动模型
type Activity struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
Date        time.Time `json:"date"`
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	Capacity    int       `json:"capacity"`
	Registered  int       `json:"registered"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Volunteer 志愿者模型
type Volunteer struct {
ID         uint      `json:"id" gorm:"primarykey"`
UserID     uint      `json:"user_id" gorm:"uniqueIndex;not null"`
User       User      `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
Name       string    `json:"name"`
Phone      string    `json:"phone"`
Email      string    `json:"email"`
Address    string    `json:"address"`
Hours      int       `json:"hours"`
Activities int       `json:"activities"`
Status     string    `json:"status"`
CreatedAt  time.Time `json:"created_at"`
UpdatedAt  time.Time `json:"updated_at"`
}

// Registration 报名记录模型
type Registration struct {
ID               uint      `json:"id" gorm:"primarykey"`
ActivityID       uint      `json:"activity_id"`
UserID           uint      `json:"user_id"` // 添加UserID字段
Name             string    `json:"name"`
Phone            string    `json:"phone"`
IDCard           string    `json:"id_card"`
Email            string    `json:"email"`
EmergencyContact string    `json:"emergency_contact"`
EmergencyPhone   string    `json:"emergency_phone"`
Status           string    `json:"status"`
CreateTime       time.Time `json:"create_time"`
UpdatedAt        time.Time `json:"updated_at"`
}

// RegistrationRequest 活动报名请求
type RegistrationRequest struct {
Name            string `json:"name" binding:"required" example:"张三"`
Phone           string `json:"phone" binding:"required" example:"13800138000"`
IDCard          string `json:"id_card" binding:"required" example:"110101199001011234"`
Email           string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
EmergencyContact string `json:"emergency_contact" binding:"required" example:"李四"`
EmergencyPhone   string `json:"emergency_phone" binding:"required" example:"13900139000"`
}
