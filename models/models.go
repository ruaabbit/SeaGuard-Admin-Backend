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
	Date        string    `json:"date"`
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
