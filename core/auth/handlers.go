package auth

import (
	"errors"
	"github.com/chfanghr/hydric/core/auth/models"
	"github.com/chfanghr/hydric/core/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

func (s *Service) LoginHandler(c *gin.Context) {
	var body models.Login

	_ = c.BindJSON(&body)
	err := body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	user, err := s.GetUserByEmail(body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, shared.ErrorMessageResponse("invalid email or password"))
		return
	}

	pwd := []byte(body.Password)
	pwdMatched := ComparePasswords(user.Password, pwd)

	if pwdMatched {
		tok := models.NewTokenFor(user)
		if err := s.InsertUserLoginToken(tok); err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
			"token": tok.Token,
			"uid":   tok.UID,
		}))
		return
	}

	c.JSON(http.StatusUnauthorized, shared.ErrorMessageResponse("invalid email or password"))
}

func (s *Service) SignupHandler(c *gin.Context) {
	var body models.Signup

	_ = c.BindJSON(&body)
	err := body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	if _, err := s.GetUserByEmail(body.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser, err := s.CreateUser(body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, shared.ErrorMessageResponse(err.Error()))
				return
			}
			c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
				"uid":   newUser.ID,
				"email": newUser.Email,
			}))
			return
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusConflict, shared.ErrorMessageResponse("user already exists."))
}

func (s *Service) IsUserExistsHandler(c *gin.Context) {
	var body models.Signup

	_ = c.BindJSON(&body)

	if _, err := s.GetUserByEmail(body.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
				"email":         body.Email,
				"alreadyExists": false,
			}))
			return
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, shared.SuccessDataResponse(gin.H{
		"email":         body.Email,
		"alreadyExists": true,
	}))
}

func (s *Service) LogoutHandler(c *gin.Context) {
	uid, _ := c.Get("uid")
	tok, _ := c.Get("token")

	if err := s.RemoveUserLoginToken(models.NewToken(uid.(uint), tok.(uuid.UUID))); err != nil {
		c.JSON(http.StatusBadRequest, shared.ErrorMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, shared.SuccessMessageResponse("logout"))
}
