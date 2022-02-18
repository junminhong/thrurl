package http

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"net/http"
	"time"
)

type shortenUrlHandler struct {
	shortenUrlUseCase domain.ShortenUrlUseCase
	shortenUrlRepo    domain.ShortenUrlRepository
}

func NewShortenUrlHandler(router *gin.Engine, shortenUrlUseCase domain.ShortenUrlUseCase, shortenUrlRepo domain.ShortenUrlRepository) {
	handler := &shortenUrlHandler{shortenUrlUseCase, shortenUrlRepo}
	router.POST("/api/v1/url", handler.ShortenUrl)
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
func (shortenUrlHandler *shortenUrlHandler) ShortenUrl(c *gin.Context) {
	request := requester.ShortenUrl{}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.ReqBindErr.Code(),
			Message:    responser.ReqBindErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	//atomicToken := requester.GetAtomicToken(c)

	response := shortenUrlHandler.shortenUrlUseCase.ShortenUrl(request)
	c.JSON(http.StatusOK, response)
}
