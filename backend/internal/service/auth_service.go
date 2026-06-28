package service

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"zblog-backend/internal/config"
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
	"zblog-backend/internal/utils"
)

type AuthService interface {
	Login(username, password string) (string, error)
	GetProfile(userID uint) (*model.User, error)
	InitAdmin() error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.WithField("username", username).Warn("login failed: user not found")
			return "", fmt.Errorf("username or password incorrect")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		logrus.WithField("username", username).Warn("login failed: password incorrect")
		return "", fmt.Errorf("username or password incorrect")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, config.Cfg.JWT.Secret, config.Cfg.JWT.Expire)
	if err != nil {
		return "", err
	}

	logrus.WithField("username", username).Info("user logged in")
	return token, nil
}

func (s *authService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) InitAdmin() error {
	count, err := s.userRepo.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username:     "admin",
		PasswordHash: string(hash),
		Nickname:     "Admin",
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(admin); err != nil {
		return err
	}

	logrus.Info("default admin user created (username: admin, password: admin123)")
	return nil
}
