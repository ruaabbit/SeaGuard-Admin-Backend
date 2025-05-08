package service

import (
	"seaguard-admin-backend/models"
	"seaguard-admin-backend/repository"
	"time"
)

// VolunteerService 志愿者服务接口
type VolunteerService interface {
GetAllVolunteers() ([]models.Volunteer, error)
CreateVolunteer(volunteer *models.Volunteer) error
UpdateVolunteer(id uint, volunteer *models.Volunteer) error
DeleteVolunteer(id uint) error
UpdateVolunteerInfo(userID uint, req *models.UpdateVolunteerInfoRequest) error
GetVolunteerInfo(userID uint) (*models.Volunteer, error)
FindByUserID(userID uint) (*models.Volunteer, error)
}

type volunteerService struct {
	repo repository.VolunteerRepository
}

// NewVolunteerService 创建志愿者服务实例
func NewVolunteerService(repo repository.VolunteerRepository) VolunteerService {
	return &volunteerService{
		repo: repo,
	}
}

// GetAllVolunteers 获取所有志愿者
func (s *volunteerService) GetAllVolunteers() ([]models.Volunteer, error) {
	return s.repo.FindAll()
}

// CreateVolunteer 创建志愿者
func (s *volunteerService) CreateVolunteer(volunteer *models.Volunteer) error {
	volunteer.Hours = 0
	volunteer.Activities = 0
	volunteer.Status = "活跃"
	volunteer.CreatedAt = time.Now()
	volunteer.UpdatedAt = time.Now()
	return s.repo.Create(volunteer)
}

// UpdateVolunteer 更新志愿者
func (s *volunteerService) UpdateVolunteer(id uint, volunteer *models.Volunteer) error {
	existingVolunteer, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	volunteer.ID = id
	volunteer.UpdatedAt = time.Now()
	volunteer.Hours = existingVolunteer.Hours
	volunteer.Activities = existingVolunteer.Activities
	return s.repo.Update(volunteer)
}

// DeleteVolunteer 删除志愿者
func (s *volunteerService) DeleteVolunteer(id uint) error {
return s.repo.Delete(id)
}

// GetVolunteerInfo 获取志愿者个人信息
func (s *volunteerService) GetVolunteerInfo(userID uint) (*models.Volunteer, error) {
    return s.repo.FindByUserID(userID)
}

// FindByUserID 根据用户ID查找志愿者
func (s *volunteerService) FindByUserID(userID uint) (*models.Volunteer, error) {
    return s.repo.FindByUserID(userID)
}

// UpdateVolunteerInfo 更新志愿者个人信息
func (s *volunteerService) UpdateVolunteerInfo(userID uint, req *models.UpdateVolunteerInfoRequest) error {
existingVolunteer, err := s.repo.FindByUserID(userID)
	if err != nil {
		return err
	}

	// 只更新允许的字段
	existingVolunteer.Name = req.Name
	existingVolunteer.Phone = req.Phone
	existingVolunteer.Email = req.Email
	existingVolunteer.Address = req.Address
	existingVolunteer.UpdatedAt = time.Now()

	return s.repo.Update(existingVolunteer)
}
