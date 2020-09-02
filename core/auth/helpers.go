package auth

import (
	"errors"
	"github.com/chfanghr/hydric/core/auth/models"
	models2 "github.com/chfanghr/hydric/core/models"
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	return bcrypt.CompareHashAndPassword(byteHash, plainPwd) == nil
}

func (s *Service) GetUserByEmail(email string) (*models2.User, error) {
	var user models2.User

	if err := s.DB.Model(&user).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) InsertUserLoginToken(t *models.Token) error {
	return s.Redis.SAdd(t.RedisKey(), t.RedisToken()).Err()
}

func (s *Service) RemoveUserLoginToken(t *models.Token) error {
	return s.Redis.SRem(t.RedisKey(), t.RedisToken()).Err()
}

func (s *Service) CreateUser(u models.Signup) (*models2.User, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	newUser := models2.User{
		Email:    u.Email,
		Password: string(pass),
	}

	if err := s.DB.Model(&models2.User{}).Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *Service) GetUserByToken(t *models.Token) (*models2.User, error) {
	redisRes := s.Redis.SIsMember(t.RedisKey(), t.RedisToken())
	isTokenExist := redisRes.Err() == nil && redisRes.Val()
	if !isTokenExist {
		return nil, errors.New("login token not exists")
	}

	var user models2.User
	if err := s.DB.Model(&models2.User{}).First(&user, t.UID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
