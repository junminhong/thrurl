package http

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"github.com/mssola/user_agent"
	"net/http"
	"time"
)

type urlHandler struct {
	urlApp domain.UrlApp
}

func NewUrlHandler(router *gin.Engine, urlApp domain.UrlApp) {
	handler := urlHandler{urlApp: urlApp}
	router.GET("/api/v1/url/check-safe", handler.checkUrlSafe)
	router.GET("/api/v1/url/record", handler.recordWhoClick)

}

func (urlHandler *urlHandler) checkUrlSafe(c *gin.Context) {
	sourceUrl := c.Query("source-url")
	if !handler.UrlLifeCheck(sourceUrl) {
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.NotFoundShortUrlErr.Code(),
			Message:    responser.NotFoundShortUrlErr.Reload("不是有效連結").Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	result, resultType := handler.SafeUrlCheck(sourceUrl)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: responser.CheckUrlSafeOk.Code(),
		Message:    responser.CheckUrlSafeOk.Message(),
		Data:       responser.CheckUrlSafe{Result: result, Type: resultType},
		TimeStamp:  time.Now(),
	})
}

func (urlHandler *urlHandler) recordWhoClick(c *gin.Context) {
	trackerID := c.Query("tracker-id")
	if trackerID == "" {
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.ReqBindErr.Code(),
			Message:    responser.ReqBindErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	ua := user_agent.New(c.GetHeader("User-Agent"))
	browserName, browserVersion := ua.Browser()
	shortUrlInfo := requester.RecordShortUrlInfo{
		ClickerIP:      c.ClientIP(),
		Browser:        browserName,
		BrowserVersion: browserVersion,
		Platform:       ua.Platform(),
		OS:             ua.OS(),
	}
	resultCode, message := urlHandler.urlApp.RecordShortUrlClickInfo(trackerID, shortUrlInfo)
	c.JSON(handler.OK, handler.Response{
		ResultCode: resultCode,
		Message:    message,
		Data:       "",
		TimeStamp:  time.Now().UTC(),
	})
}
