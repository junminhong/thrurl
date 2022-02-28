package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"net/http"
	"time"
)

func CheckAtomicTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		atomicToken := requester.GetAtomicToken(c)
		if atomicToken == "" {
			c.JSON(http.StatusOK, responser.Response{
				ResultCode: responser.NotFoundAtomicTokenErr.Code(),
				Message:    responser.NotFoundAtomicTokenErr.Message(),
				Data:       "",
				TimeStamp:  time.Now(),
			})
			return
		}
		c.Next()
	}
}
