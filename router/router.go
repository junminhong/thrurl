package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	v1 "github.com/junminhong/thrurl/api/v1"
	"github.com/junminhong/thrurl/pkg/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load env file")
	}
}

func middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokens := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")
		response := handler.Response{
			ResultCode: handler.TokenError1,
			Message:    handler.ErrorFlag[handler.TokenError1],
			Data:       "",
		}
		if len(tokens) != 2 {
			response.TimeStamp = time.Now().UTC()
			c.AbortWithStatusJSON(handler.OK, response)
			return
		}
		if tokens[1] == "" {
			response.TimeStamp = time.Now().UTC()
			c.AbortWithStatusJSON(handler.OK, response)
			return
		}
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
		apiRouter.POST("/short-url", v1.ShortUrl)
	}
	needTokenRouter := router.Group("api/v1").Use(checkToken())
	{
		needTokenRouter.GET("/url-list", v1.AllUrlList)
		needTokenRouter.GET("/url-paginate", v1.AllUrlPaginate)
	}
	indexRouting := router.Group("/")
	{
		indexRouting.GET("/:short-id", v1.Test)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9020"
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)))
	router.Run(":" + port)
}
