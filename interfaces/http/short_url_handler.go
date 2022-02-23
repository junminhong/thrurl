package http

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"log"
	"net/http"
	"time"
)

type shortUrlHandler struct {
	shortUrlApp  domain.ShortUrlApp
	shortUrlRepo domain.ShortUrlRepository
}

func NewShortenUrlHandler(router *gin.Engine, shortUrlApp domain.ShortUrlApp, shortenUrlRepo domain.ShortUrlRepository) {
	handler := &shortUrlHandler{shortUrlApp: shortUrlApp, shortUrlRepo: shortenUrlRepo}
	router.GET("/:tracker-id", handler.redirect)
	router.POST("/api/v1/short-url", handler.shortenUrl)
	router.GET("/api/v1/short-url", handler.getSourceUrl)
	router.PUT("/api/v1/short-url", handler.editShortUrl)
	router.GET("/api/v1/short-url/list", handler.getShortUrlList)
}

// ShortenUrl
// @Summary 縮短網址
// @Description
// @Tags url
// @version 1.0
// @Accept application/json
// @produce application/json
// @param data body requester.Register true "請求資料"
// @Success 1004 {object} responser.Response "帳戶註冊成功"
// @failure 1000 {object} responser.Response "request格式錯誤"
// @failure 1002 {object} responser.Response "信箱已經存在"
// @failure 1003 {object} responser.Response "帳戶註冊失敗"
// @Router /url [post]
func (shortUrlHandler *shortUrlHandler) shortenUrl(c *gin.Context) {
	atomicToken := requester.GetAtomicToken(c)
	request := requester.ShortenUrl{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.ReqBindErr.Code(),
			Message:    responser.ReqBindErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	resultCode, message, data := shortUrlHandler.shortUrlApp.ShortenUrl(request.SourceUrl, atomicToken)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: resultCode,
		Message:    message,
		Data:       data,
		TimeStamp:  time.Now(),
	})
}

func (shortUrlHandler *shortUrlHandler) getSourceUrl(c *gin.Context) {
}

func (shortUrlHandler *shortUrlHandler) editShortUrl(c *gin.Context) {
	atomicToken := requester.GetAtomicToken(c)
	if atomicToken == "" {
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.NotFoundAtomicTokenErr.Code(),
			Message:    responser.NotFoundAtomicTokenErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	request := requester.EditShortUrl{}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.ReqBindErr.Code(),
			Message:    responser.ReqBindErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	resultCode, message := shortUrlHandler.shortUrlApp.EditShortUrl(request, atomicToken)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: resultCode,
		Message:    message,
		Data:       "",
		TimeStamp:  time.Now(),
	})
}

func (shortUrlHandler *shortUrlHandler) getShortUrlList(c *gin.Context) {
	// page size、offset、limit
}

func (shortUrlHandler *shortUrlHandler) redirect(c *gin.Context) {
	sourceUrl := shortUrlHandler.shortUrlApp.GetSourceUrl(c.Param("tracker-id"))
	c.Redirect(http.StatusMovedPermanently, sourceUrl)
}
