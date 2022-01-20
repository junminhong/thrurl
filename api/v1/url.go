package v1

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/junminhong/thrurl/db/blot"
	"github.com/junminhong/thrurl/pkg/handler"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var blotDB = blot.Setup()

type shortUrlReq struct {
	SourceUrl  string `json:"source_url" binding:"required"`
	Expired    string `json:"expired"`
	SourceUrlB string `json:"source_url_b"`
	WhoClick   bool   `json:"who_click"`
}

type data struct {
	ShortUrl string `json:"short_url"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load env file")
	}
}

func Short(ctx *gin.Context) {
	request := &shortUrlReq{}
	err := ctx.BindJSON(request)
	response := handler.Response{}
	if err != nil {
		response = handler.Response{
			ResultCode: handler.BadRequest,
			Message:    handler.ResponseFlag[handler.BadRequest],
			Data:       &data{ShortUrl: ""},
			TimeStamp:  time.Now().UTC(),
		}
		ctx.JSON(response.ResultCode, response)
		log.Println(err)
		return
	}
	auth := ctx.Request.Header.Get("Authorization")
	if auth != "" {
		// 代表有登入
		token := strings.Split(auth, "Bearer ")[1]
		shorUrlLoginHandler(token, request)
	} else {
		response, err = shorUrlNotLoginHandler(request)
		if err != nil {
			log.Println(err)
		}
		ctx.JSON(response.ResultCode, response)
	}
}

func shorUrlLoginHandler(token string, request *shortUrlReq) {
	// 處理過期、AB測試、點擊成效
}

func shorUrlNotLoginHandler(request *shortUrlReq) (handler.Response, error) {
	shortUrl := uuid.New()
	err := blotDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UrlBucket"))
		err := bucket.Put([]byte(strings.Split(shortUrl.String(), "-")[0]), []byte(request.SourceUrl))
		return err
	})
	if err != nil {
		response := handler.Response{
			ResultCode: handler.BadRequest,
			Message:    handler.ResponseFlag[handler.BadRequest],
			Data:       data{ShortUrl: ""},
			TimeStamp:  time.Now().UTC(),
		}
		return response, err
	}
	response := handler.Response{
		ResultCode: handler.OK,
		Message:    handler.ResponseFlag[handler.OK],
		Data:       data{ShortUrl: os.Getenv("HOST_NAME") + "/" + strings.Split(shortUrl.String(), "-")[0]},
		TimeStamp:  time.Now().UTC(),
	}
	return response, err
}

func Test(ctx *gin.Context) {
	shortUrl := ""
	_ = blotDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UrlBucket"))
		shortUrl = string(bucket.Get([]byte(ctx.Param("short-url"))))
		return nil
	})
	ctx.Redirect(http.StatusMovedPermanently, shortUrl)
}
