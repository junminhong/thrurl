package postgresql

import (
	"fmt"
	"github.com/junminhong/thrurl/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBINFO struct {
	dbUser     string
	dbPassword string
	dbHost     string
	dbName     string
	dbPort     string
}

func (dbInfo DBINFO) setupDBInfo() *DBINFO {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load env file")
	}
	dbInfo.dbUser = os.Getenv("DB_USER")
	dbInfo.dbPassword = os.Getenv("DB_PASSWORD")
	dbInfo.dbHost = os.Getenv("DB_HOST")
	dbInfo.dbName = os.Getenv("DB_NAME")
	dbInfo.dbPort = os.Getenv("DB_PORT")
	return &dbInfo
}

func Setup() *gorm.DB {
	dbInfo := &DBINFO{}
	dbInfo = dbInfo.setupDBInfo()
	// sslmode=disable
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  TimeZone=Asia/Taipei",
		dbInfo.dbHost, dbInfo.dbUser, dbInfo.dbPassword, dbInfo.dbName, dbInfo.dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("Failed to connect DB")
	}
	return db
}
func MigrateDB() {
	postgresDB := Setup()
	if postgresDB.Migrator().HasTable(&model.ShortUrl{}) && postgresDB.Migrator().HasTable(&model.ShortUrlInfo{}) {
		return
	}
	err := postgresDB.AutoMigrate(&model.ShortUrl{}, &model.ShortUrlInfo{})
	if err != nil {
		log.Println(err.Error())
	}
}
