package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/junminhong/thrurl/api/v1"
	"net/http"
	"os"
)

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func Setup() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")
	router.Static("/static", "./static")
	router.Use(middleware())
	apiRouter := router.Group("api/v1")
	{
		apiRouter.POST("/short-url", v1.Short)
	}
	indexRouting := router.Group("/")
	{
		indexRouting.GET("", getIndex)
		indexRouting.GET("/login", getLogin)
		indexRouting.GET("/:short-url", v1.Test)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9020"
	}
	router.Run(":" + port)
}

func getIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Thrurl - 短網址專家",
	})
}
func getLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Thrurl - 短網址專家",
	})
}
