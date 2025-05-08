package service

import (
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/repository"
	"time"
)

// ActivityService 活动服务接口
type ActivityService interface {
GetAllActivities() ([]models.Activity, error)
GetAvailableActivities() ([]models.Activity, error)
CreateActivity(activity *models.Activity) error
UpdateActivity(id uint, activity *models.Activity) error
DeleteActivity(id uint) error
}

type activityService struct {
	repo repository.ActivityRepository
}

// NewActivityService 创建活动服务实例
func NewActivityService(repo repository.ActivityRepository) ActivityService {
	return &activityService{
		repo: repo,
	}
}

// GetAllActivities 获取所有活动
func (s *activityService) GetAllActivities() ([]models.Activity, error) {
return s.repo.FindAll()
}

// GetAvailableActivities 获取可报名活动
func (s *activityService) GetAvailableActivities() ([]models.Activity, error) {
    activities, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }
    
    now := time.Now()
    var availableActivities []models.Activity
    for _, activity := range activities {
        // 检查活动状态、日期和容量
        if activity.Status == "报名中" && 
           activity.Date.After(now) && 
           activity.Registered < activity.Capacity {
            availableActivities = append(availableActivities, activity)
        }
    }
    return availableActivities, nil
}

// CreateActivity 创建活动
func (s *activityService) CreateActivity(activity *models.Activity) error {
	activity.Status = "报名中"
	activity.Registered = 0
	activity.CreatedAt = time.Now()
	activity.UpdatedAt = time.Now()
	return s.repo.Create(activity)
}

// UpdateActivity 更新活动
func (s *activityService) UpdateActivity(id uint, activity *models.Activity) error {
	existingActivity, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	activity.ID = id
	activity.UpdatedAt = time.Now()
	activity.Registered = existingActivity.Registered
	return s.repo.Update(activity)
}

// DeleteActivity 删除活动
func (s *activityService) DeleteActivity(id uint) error {
	return s.repo.Delete(id)
}
