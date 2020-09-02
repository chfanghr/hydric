package auth

import (
	"github.com/chfanghr/hydric/core/config"
	models2 "github.com/chfanghr/hydric/core/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"log"
)

type Service struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (s *Service) SetupDB() error {
	return s.DB.AutoMigrate(&models2.User{})
}

func NewService(DB *gorm.DB, redis *redis.Client) (*Service, error) {
	s := &Service{DB: DB, Redis: redis}
	return s, s.SetupDB()
}

func FromConfig(c *config.Parsed) (*Service, error) {
	s := &Service{
		c.DBConn,
		c.RedisClient,
	}
	return s, s.SetupDB()
}

func MustFromConfig(c *config.Parsed) *Service {
	s, err := FromConfig(c)
	if err != nil {
		log.Panicln(err)
	}
	return s
}

func (s *Service) SetupDefaultAuthAPI(g gin.IRouter) {
	g.POST("/login", s.LoginHandler)
	g.POST("/signup", s.SignupHandler)
	g.POST("/logout", s.MakeAuthRequiredMiddleware(), s.LogoutHandler)
	g.GET("/logout", s.MakeAuthRequiredMiddleware(), s.LogoutHandler)
}
