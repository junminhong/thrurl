package http

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"github.com/mssola/user_agent"
	"net"
	"net/http"
	"strings"
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

// checkUrlSafe
// @Summary 檢查網址安全
// @Description
// @Tags url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param  source-url  query string  true  "source-url"
// @Success 1007 {object} responser.CheckUrlSafe "檢查完成"
// @failure 1005 {object} string "不是有效連結"
// @Router /api/v1/url/check-safe [get]
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

// recordWhoClick
// @Summary 記錄點擊成效
// @Description
// @Tags url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param tracker-id query string  true "tracker-id"
// @Success 1009 {object} responser.CheckUrlSafe "記錄完成"
// @failure 1000 {object} string "請依照API文件發起請求"
// @failure 1008 {object} string "記錄失敗"
// @Router /api/v1/url/record [get]
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
	clientIP := net.ParseIP(c.ClientIP()).To4()
	clientCountry := c.GetHeader("CF-IPCountry")
	if strings.Compare(c.GetHeader("CF-IPCountry"), "XX") == 0 {
		clientCountry = ""
	}
	shortUrlInfo := requester.RecordShortUrlInfo{
		ClickerIP:      clientIP.String(),
		ClickerCountry: clientCountry,
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
