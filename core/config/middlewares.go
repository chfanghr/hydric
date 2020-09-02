package config

import (
	"github.com/chfanghr/hydric/core/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

const RedisUnderMaintenanceKey = "is_under_maintenance"

func (p *Parsed) MakeCheckMaintenanceStatusMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		existRes := p.RedisClient.Exists(RedisUnderMaintenanceKey)
		if err := existRes.Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorMessageResponse(err.Error()))
			return
		}
		keyExists := existRes.Val() == 1
		if keyExists {
			getRes := p.RedisClient.Get(RedisUnderMaintenanceKey)
			if err := getRes.Err(); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, shared.ErrorMessageResponse(err.Error()))
				return
			}
			val := getRes.Val()
			if val == "1" || val == "true" {
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}
		}
		c.Next()
	}
}
