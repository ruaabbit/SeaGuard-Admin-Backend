package repository

import (
	"seaguard-admin-backend/config"
	"seaguard-admin-backend/models"
)

// VolunteerRepository 志愿者仓储接口
type VolunteerRepository interface {
FindAll() ([]models.Volunteer, error)
Create(volunteer *models.Volunteer) error
FindByID(id uint) (*models.Volunteer, error)
FindByUserID(userID uint) (*models.Volunteer, error)
Update(volunteer *models.Volunteer) error
Delete(id uint) error
}

type volunteerRepository struct{}

// NewVolunteerRepository 创建志愿者仓储实例
func NewVolunteerRepository() VolunteerRepository {
	return &volunteerRepository{}
}

// FindAll 获取所有志愿者
func (r *volunteerRepository) FindAll() ([]models.Volunteer, error) {
var volunteers []models.Volunteer
err := config.DB.Preload("User").Find(&volunteers).Error
	return volunteers, err
}

// Create 创建志愿者
func (r *volunteerRepository) Create(volunteer *models.Volunteer) error {
	return config.DB.Create(volunteer).Error
}

// FindByID 根据ID查找志愿者
func (r *volunteerRepository) FindByID(id uint) (*models.Volunteer, error) {
var volunteer models.Volunteer
err := config.DB.Preload("User").First(&volunteer, id).Error
return &volunteer, err
}

// FindByUserID 根据UserID查找志愿者
func (r *volunteerRepository) FindByUserID(userID uint) (*models.Volunteer, error) {
var volunteer models.Volunteer
err := config.DB.Preload("User").Where("user_id = ?", userID).First(&volunteer).Error
return &volunteer, err
}

// Update 更新志愿者
func (r *volunteerRepository) Update(volunteer *models.Volunteer) error {
	return config.DB.Save(volunteer).Error
}

// Delete 删除志愿者
func (r *volunteerRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Volunteer{}, id).Error
}
