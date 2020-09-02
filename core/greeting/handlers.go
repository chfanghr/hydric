package greeting

import (
	"fmt"
	"github.com/chfanghr/hydric/core/models"
	"github.com/chfanghr/hydric/core/shared"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GreetHandler(c *gin.Context) {
	user, isLogin := c.Get("user")
	name := c.Param("name")
	if len(name) > 0 {
		if isLogin {
			c.JSON(http.StatusOK, shared.SuccessMessageResponse(Greet(fmt.Sprintf("%s(%s)", user.(*models.User).Email, name))))
			return
		}
		c.JSON(http.StatusOK, shared.SuccessMessageResponse(Greet(name)))
		return
	}

	c.JSON(http.StatusBadRequest, shared.SuccessMessageResponse("tell me your name please"))
}
