package v1

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/junminhong/thrurl/db/blot"
	"log"
	"net/http"
	"strings"
)

var blotDB = blot.Setup()

type shortReq struct {
	Url     string `json:"url"`
	Expired string `json:"expired"`
}

func Short(ctx *gin.Context) {
	request := &shortReq{}
	err := ctx.BindJSON(request)
	log.Println(request.Url)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "請依照API文件請求",
		})
		return
	}
	shortUrl := uuid.New()
	err = blotDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UrlBucket"))
		err = bucket.Put([]byte(strings.Split(shortUrl.String(), "-")[0]), []byte(request.Url))
		return err
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "處理失敗，請洽詢後端",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "處理完成",
		"shortUrl": " https://thrurl.herokuapp.com/" + strings.Split(shortUrl.String(), "-")[0],
	})
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
