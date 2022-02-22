package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	deliver "github.com/junminhong/thrurl/api/v1/delivery/http"
	"github.com/junminhong/thrurl/api/v1/delivery/http/middleware"
	"github.com/junminhong/thrurl/api/v1/repository"
	"github.com/junminhong/thrurl/api/v1/usecase"
	_ "github.com/junminhong/thrurl/docs"
	"github.com/junminhong/thrurl/domain"
	"github.com/spf13/viper"
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
	err := postgresDB.db.AutoMigrate(&domain.ShortenUrl{}, &domain.ShortenUrlInfo{})
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
	shortenUrlRepo := repository.NewShortenUrlRepo(db, redis)
	shortenUrlCase := usecase.NewShortenUrlUseCase(shortenUrlRepo)
	deliver.NewShortenUrlHandler(router, grpcClient, shortenUrlCase, shortenUrlRepo)
}

func setUpGrpcClient() *grpc.ClientConn {
	conn, err := grpc.Dial("127.0.0.1:9205", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return conn
}

func main() {
	router := gin.Default()
	router.Use(middleware.Middleware())
	db := setUpDB()
	redisClient := setUpRedis()
	grpcClient := setUpGrpcClient()
	setUpDomain(router, grpcClient, db.db, redisClient)
	//db.migrationDB()
	router.Run("127.0.0.1:9220")
}
