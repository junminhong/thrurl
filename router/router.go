package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/junminhong/thrurl/api/v1"
	"net/http"
	"os"
)

func Setup() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")
	apiRouter := router.Group("api/v1")
	{
		apiRouter.POST("/short-url", v1.Short)
	}
	indexRouting := router.Group("/")
	{
		indexRouting.GET("", getIndex)
		indexRouting.GET("/:short-url", v1.Test)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9020"
	}
	router.Run(":" + port)
}

func getIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}
