package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/junminhong/thrurl/application"
	_ "github.com/junminhong/thrurl/docs"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/infrastructure/repo"
	deliver "github.com/junminhong/thrurl/interfaces/http"
	"github.com/junminhong/thrurl/interfaces/http/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

// @title           Thrurl API
// @version         1.0
// @description     一個簡單易用且強大的縮網址服務
// @termsOfService  http://swagger.io/terms/

// @contact.name   junmin.hong
// @contact.url    https://github.com/junminhong
// @contact.email  junminhong1110@gmail.com

// @license.name  MIT
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:9020
// @BasePath  /api/v1
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
	}
}

type postgresDB struct {
	db *gorm.DB
}

func setUpDB() *postgresDB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Taipei",
		viper.GetString("APP.DB_HOST"),
		viper.GetString("APP.DB_USERNAME"),
		viper.GetString("APP.DB_PASSWORD"),
		viper.GetString("APP.DB_DATABASE"),
		viper.GetString("APP.DB_PORT"),
	)
	log.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("Failed to connect DB")
	}
	return &postgresDB{db: db}
}
func (postgresDB *postgresDB) migrationDB() {
	err := postgresDB.db.AutoMigrate(&domain.ShortUrl{}, &domain.ShortUrlInfo{}, &domain.ClickInfo{})
	if err != nil {
		log.Println(err.Error())
	}
}
func setUpRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("APP.REDIS_HOST") + ":" + viper.GetString("APP.REDIS_PORT"),
		Password: viper.GetString("APP.REDIS_PASSWORD"), // no password set
		DB:       0,                                     // use default DB
	})
	return client
}
func setUpDomain(router *gin.Engine, grpcClient *grpc.ClientConn, db *gorm.DB, redis *redis.Client) {
	shortUrlRepo := repo.NewShortenUrlRepo(db, redis, grpcClient)
	shortUrlApp := application.NewShortenUrlUseCase(shortUrlRepo)
	deliver.NewShortenUrlHandler(router, shortUrlApp)
	urlRepo := repo.NewUrlRepo(db)
	urlApp := application.NewUrlApp(urlRepo)
	deliver.NewUrlHandler(router, urlApp)
}

func setUpGrpcClient() *grpc.ClientConn {
	conn, err := grpc.Dial(viper.GetString("APP.GRPC_HOST")+":"+viper.GetString("APP.GRPC_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return conn
}

func setUpRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Middleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)))
	return router
}

func main() {
	db := setUpDB()
	go db.migrationDB()
	redisClient := setUpRedis()
	grpcClient := setUpGrpcClient()
	router := setUpRouter()
	setUpDomain(router, grpcClient, db.db, redisClient)
	router.Run(viper.GetString("APP.HOST") + ":" + viper.GetString("APP.PORT"))
}
