package repository

import (
	"seaguard-admin-backend/config"
	"seaguard-admin-backend/models"
)

// ActivityRepository 活动仓储接口
type ActivityRepository interface {
	FindAll() ([]models.Activity, error)
	Create(activity *models.Activity) error
	FindByID(id uint) (*models.Activity, error)
	Update(activity *models.Activity) error
	Delete(id uint) error
}

type activityRepository struct{}

// NewActivityRepository 创建活动仓储实例
func NewActivityRepository() ActivityRepository {
	return &activityRepository{}
}

// FindAll 获取所有活动
func (r *activityRepository) FindAll() ([]models.Activity, error) {
	var activities []models.Activity
	err := config.DB.Find(&activities).Error
	return activities, err
}

// Create 创建活动
func (r *activityRepository) Create(activity *models.Activity) error {
	return config.DB.Create(activity).Error
}

// FindByID 根据ID查找活动
func (r *activityRepository) FindByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	err := config.DB.First(&activity, id).Error
	return &activity, err
}

// Update 更新活动
func (r *activityRepository) Update(activity *models.Activity) error {
	return config.DB.Save(activity).Error
}

// Delete 删除活动
func (r *activityRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Activity{}, id).Error
}
