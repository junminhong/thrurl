package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	v1 "github.com/junminhong/thrurl/api/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
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

func Setup() {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")
	router.Static("/static", "./static")
	router.Use(middleware())
	apiRouter := router.Group("api/v1")
	{
		apiRouter.POST("/short-url", v1.ShortUrl)
	}
	indexRouting := router.Group("/")
	{
		indexRouting.GET("/:short-url", v1.Test)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "9020"
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)))
	router.Run(":" + port)
}
