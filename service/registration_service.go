package service

import (
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/repository"
	"time"
)

// RegistrationService 报名服务接口
type RegistrationService interface {
	GetActivityRegistrations(activityID uint) ([]models.Registration, error)
	UpdateRegistrationStatus(id uint, status string) error
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
