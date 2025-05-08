package repository

import (
	"seaguard-admin-backend/config"
	"seaguard-admin-backend/models"
)

// RegistrationRepository 报名记录仓储接口
type RegistrationRepository interface {
FindByActivityID(activityID uint) ([]models.Registration, error)
FindByID(id uint) (*models.Registration, error)
Create(registration *models.Registration) error
Update(registration *models.Registration) error
UpdateStatus(id uint, status string) error
FindByUserAndActivity(userID, activityID uint) (*models.Registration, error)
CheckDuplicateRegistration(userID, activityID uint) (bool, error)
}

type registrationRepository struct{}

// NewRegistrationRepository 创建报名记录仓储实例
func NewRegistrationRepository() RegistrationRepository {
	return &registrationRepository{}
}

// FindByActivityID 获取活动的所有报名记录
func (r *registrationRepository) FindByActivityID(activityID uint) ([]models.Registration, error) {
	var registrations []models.Registration
	err := config.DB.Where("activity_id = ?", activityID).Find(&registrations).Error
	return registrations, err
}

// FindByID 根据ID查找报名记录
func (r *registrationRepository) FindByID(id uint) (*models.Registration, error) {
	var registration models.Registration
	err := config.DB.First(&registration, id).Error
	return &registration, err
}

// Create 创建报名记录
func (r *registrationRepository) Create(registration *models.Registration) error {
	return config.DB.Create(registration).Error
}

// Update 更新报名记录
func (r *registrationRepository) Update(registration *models.Registration) error {
	return config.DB.Save(registration).Error
}

// UpdateStatus 更新报名状态
func (r *registrationRepository) UpdateStatus(id uint, status string) error {
return config.DB.Model(&models.Registration{}).Where("id = ?", id).Update("status", status).Error
}

// FindByUserAndActivity 查找用户在某个活动的报名记录
func (r *registrationRepository) FindByUserAndActivity(userID, activityID uint) (*models.Registration, error) {
    var registration models.Registration
    err := config.DB.Where("user_id = ? AND activity_id = ?", userID, activityID).First(&registration).Error
    return &registration, err
}

// CheckDuplicateRegistration 检查是否重复报名
func (r *registrationRepository) CheckDuplicateRegistration(userID, activityID uint) (bool, error) {
    var count int64
    err := config.DB.Model(&models.Registration{}).
        Where("user_id = ? AND activity_id = ?", userID, activityID).
        Count(&count).Error
    return count > 0, err
}
