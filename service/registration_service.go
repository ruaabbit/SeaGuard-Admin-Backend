package service

import (
"errors"
"seaguard-admin-backend/config"
"seaguard-admin-backend/models"
"seaguard-admin-backend/repository"
"time"
)

// RegistrationService 报名服务接口
type RegistrationService interface {
GetActivityRegistrations(activityID uint) ([]models.Registration, error)
UpdateRegistrationStatus(id uint, status string) error
CreateRegistration(userID uint, registration *models.Registration) error
GetUserRegistration(userID, activityID uint) (*models.Registration, error)
}

type registrationService struct {
	regRepo repository.RegistrationRepository
	actRepo repository.ActivityRepository
}

// NewRegistrationService 创建报名服务实例
func NewRegistrationService(
	regRepo repository.RegistrationRepository,
	actRepo repository.ActivityRepository,
) RegistrationService {
	return &registrationService{
		regRepo: regRepo,
		actRepo: actRepo,
	}
}

// GetActivityRegistrations 获取活动的所有报名记录
func (s *registrationService) GetActivityRegistrations(activityID uint) ([]models.Registration, error) {
	return s.regRepo.FindByActivityID(activityID)
}

// UpdateRegistrationStatus 更新报名状态
func (s *registrationService) UpdateRegistrationStatus(id uint, status string) error {
	registration, err := s.regRepo.FindByID(id)
	if err != nil {
		return err
	}

	oldStatus := registration.Status
	registration.Status = status
	registration.UpdatedAt = time.Now()

	err = s.regRepo.Update(registration)
	if err != nil {
		return err
	}

	// 更新活动报名人数
	if status == "已通过" && oldStatus != "已通过" {
		activity, err := s.actRepo.FindByID(registration.ActivityID)
		if err != nil {
			return err
		}

		activity.Registered++
		err = s.actRepo.Update(activity)
		if err != nil {
			return err
		}
	} else if oldStatus == "已通过" && status != "已通过" {
		activity, err := s.actRepo.FindByID(registration.ActivityID)
		if err != nil {
			return err
		}

		activity.Registered--
		err = s.actRepo.Update(activity)
		if err != nil {
			return err
		}
	}

return nil
}

// CreateRegistration 创建报名记录
func (s *registrationService) CreateRegistration(userID uint, registration *models.Registration) error {
    // 检查活动是否存在及可报名
    activity, err := s.actRepo.FindByID(registration.ActivityID)
    if err != nil {
        return err
    }
    
    if activity.Status != "进行中" {
        return errors.New("活动不在报名阶段")
    }
    
    if activity.Registered >= activity.Capacity {
        return errors.New("活动名额已满")
    }

    // 检查是否重复报名
    isDuplicate, err := s.regRepo.CheckDuplicateRegistration(userID, registration.ActivityID)
    if err != nil {
        return err
    }
    if isDuplicate {
        return errors.New("already registered")
    }

    // 设置报名记录属性
    registration.UserID = userID
    registration.CreateTime = time.Now()
    registration.UpdatedAt = time.Now()
    registration.Status = "pending"

    // 开启事务
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 创建报名记录
    if err := s.regRepo.Create(registration); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

// GetUserRegistration 获取用户在某个活动的报名记录
func (s *registrationService) GetUserRegistration(userID, activityID uint) (*models.Registration, error) {
    return s.regRepo.FindByUserAndActivity(userID, activityID)
}
