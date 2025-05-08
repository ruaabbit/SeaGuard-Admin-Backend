package service

import (
"errors"
"golang.org/x/crypto/bcrypt"
"seaguard-admin-backend/models"
"seaguard-admin-backend/repository"
"seaguard-admin-backend/utils"
"seaguard-admin-backend/config"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(req *models.RegisterRequest) error {
    // 检查用户名是否已存在
    existingUser, _ := s.userRepo.FindByUsername(req.Username)
    if existingUser != nil {
        return errors.New("用户名已存在")
    }

    // 密码加密
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // 创建新用户
    user := &models.User{
        Username: req.Username,
        Password: string(hashedPassword),
        Role:     req.Role,
        Status:   "active",
    }

    // 开启事务
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 创建用户账号
    if err := tx.Create(user).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 如果是志愿者角色，创建志愿者信息
    if req.Role == "volunteer" {
        volunteer := &models.Volunteer{
            UserID:     user.ID,  // 设置关联的用户ID
            Name:       req.Name,
            Phone:      req.Phone,
            Email:      req.Email,
            Address:    req.Address,
            Hours:      0,
            Activities: 0,
            Status:     "活跃",
        }
        
        if err := tx.Create(volunteer).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit().Error
}

func (s *UserService) Login(username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		return nil, "", errors.New("用户账号已被禁用")
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, "", errors.New("生成token失败")
	}

	return user, token, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) DeleteUser(id uint) error {
    // 开启事务
    tx := config.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 查询用户
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        tx.Rollback()
        return err
    }

    // 如果是志愿者，先删除志愿者信息（由于设置了CASCADE，这步可以省略）
    if user.Role == "volunteer" {
        if err := tx.Where("user_id = ?", id).Delete(&models.Volunteer{}).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    // 删除用户
    if err := s.userRepo.Delete(id); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *UserService) UpdateStatus(userID uint, status string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Status = status
	return s.userRepo.Update(user)
}
