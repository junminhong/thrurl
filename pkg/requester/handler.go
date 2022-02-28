package requester

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetAtomicToken(c *gin.Context) (atomicToken string) {
	token := c.Request.Header.Get("Authorization")
	tokens := strings.Split(token, " ")
	if len(tokens) != 2 {
		return ""
	}
	if strings.Compare(tokens[0], "Bearer") != 0 {
		return ""
	}
	atomicToken = tokens[1]
	return atomicToken
}
