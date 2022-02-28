package http

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/interfaces/http/middleware"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type shortUrlHandler struct {
	shortUrlApp domain.ShortUrlApp
}

func NewShortenUrlHandler(router *gin.Engine, shortUrlApp domain.ShortUrlApp) {
	handler := &shortUrlHandler{shortUrlApp: shortUrlApp}
	router.GET("/:tracker-id", handler.redirect)
	router.Static("/api/v1/image", "./static")
	router.POST("/api/v1/short-url", handler.shortenUrl)
	router.GET("/api/v1/short-url/redirect", handler.getRedirectUrl)
	needAtomicToken := router.Group("/api/v1/short-url").Use(middleware.CheckAtomicTokenMiddleware())
	needAtomicToken.GET("", handler.getShortUrl)
	needAtomicToken.PUT("", handler.editShortUrl)
	needAtomicToken.GET("/list", handler.getShortUrlList)
	needAtomicToken.GET("/click-info", handler.getShortUrlClickInfo)
}

// shortenUrl
// @Summary 縮短網址
// @Description
// @Tags short-url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param Authorization header string false "Atomic Token" default(Bearer )
// @param data body requester.ShortenUrl true "請求資料"
// @Success 1002 {object} responser.Response "短連結生成成功"
// @failure 1000 {object} string "請依照API文件進行請求"
// @failure 1001 {object} string "短連結生成失敗"
// @failure 1003 {object} string "無效連結"
// @Router /api/v1/short-url [post]
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

// getShortUrl
// @Summary 取得短連結的訊息
// @Description
// @Tags short-url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param Authorization header string true "Atomic Token" default(Bearer <請在這邊輸入Atomic Token>)
// @Param  tracker-id  query string  true  "tracker-id"
// @Success 1002 {object} responser.Response "短連結生成成功"
// @failure 1000 {object} string "請依照API文件進行請求"
// @failure 1001 {object} string "短連結生成失敗"
// @failure 1003 {object} string "無效連結"
// @Router /api/v1/short-url [get]
func (shortUrlHandler *shortUrlHandler) getShortUrl(c *gin.Context) {
	atomicToken := requester.GetAtomicToken(c)
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
	resultCode, message, data := shortUrlHandler.shortUrlApp.GetShortUrl(trackerID, atomicToken)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: resultCode,
		Message:    message,
		Data:       data,
		TimeStamp:  time.Now(),
	})
}

// editShortUrl
// @Summary 編輯短連結
// @Description
// @Tags short-url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param Authorization header string true "Atomic Token" default(Bearer <請在這邊輸入Atomic Token>)
// @param data body requester.EditShortUrl true "請求資料"
// @Success 1002 {object} responser.Response "短連結保存成功"
// @failure 1000 {object} string "請依照API文件進行請求"
// @failure 1001 {object} string "短連結保存失敗"
// @failure 1003 {object} string "無效連結"
// @failure 1004 {object} string "你沒有權限發起該請求"
// @failure 1005 {object} string "找不到該連結訊息"
// @Router /api/v1/short-url [put]
func (shortUrlHandler *shortUrlHandler) editShortUrl(c *gin.Context) {
	atomicToken := requester.GetAtomicToken(c)
	request := requester.EditShortUrl{}
	if err := c.ShouldBindJSON(&request); err != nil {
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

// getShortUrlList
// @Summary 取得短連結的列表
// @Description
// @Tags short-url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param Authorization header string true "Atomic Token" default(Bearer <請在這邊輸入Atomic Token>)
// @Success 1006 {object} responser.Response "取得短連結列表"
// @failure 1004 {object} string "你沒有權限發起該請求"
// @Router /api/v1/short-url/list [get]
func (shortUrlHandler *shortUrlHandler) getShortUrlList(c *gin.Context) {
	// page size、offset、limit
	// 預設一頁10個
	atomicToken := requester.GetAtomicToken(c)
	limit := 10
	if value := c.Query("limit"); value != "" && strings.Compare(value, "0") != 0 {
		limit, _ = strconv.Atoi(value)
	}
	offset := 0
	if value := c.Query("offset"); value != "" {
		offset, _ = strconv.Atoi(value)
	}
	resultCode, message, data, page := shortUrlHandler.shortUrlApp.GetShortUrlList(limit, offset, atomicToken)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: resultCode,
		Message:    message,
		Data:       responser.ShortUrlLists{ShortUrlList: data, Page: page},
		TimeStamp:  time.Now(),
	})
}

func (shortUrlHandler *shortUrlHandler) redirect(c *gin.Context) {
	//sourceUrl := shortUrlHandler.shortUrlApp.GetSourceUrl(c.Param("tracker-id"))
	//c.Redirect(http.StatusMovedPermanently, sourceUrl)
}

// getRedirectUrl
// @Summary 取得轉址要去的目標
// @Description
// @Tags short-url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param  tracker-id  query string  true  "tracker-id"
// @Success 1010 {object} responser.Response "成功取得原始連結"
// @Router /api/v1/short-url/redirect [get]
func (shortUrlHandler *shortUrlHandler) getRedirectUrl(c *gin.Context) {
	resultCode, message, sourceUrl := shortUrlHandler.shortUrlApp.GetSourceUrl(c.Query("tracker-id"))
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: resultCode,
		Message:    message,
		Data: responser.GetSourceUrl{
			SourceUrl: sourceUrl,
		},
		TimeStamp: time.Now(),
	})
}

// getShortUrlClickInfo
// @Summary 取得短連結的點擊成效
// @Description
// @Tags short-url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Param Authorization header string true "Atomic Token" default(Bearer <請在這邊輸入Atomic Token>)
// @Param  tracker-id  query string  true  "tracker-id"
// @Success 1011 {object} responser.Response "成功取得點擊成效"
// @failure 1004 {object} string "你沒有權限發起該請求"
// @failure 1005 {object} string "找不到該短連結的點擊成效"
// @Router /api/v1/short-url/click-info [get]
func (shortUrlHandler *shortUrlHandler) getShortUrlClickInfo(c *gin.Context) {
	atomicToken := requester.GetAtomicToken(c)
	resultCode, message, data := shortUrlHandler.shortUrlApp.GetShortUrlClickInfo(c.Query("tracker-id"), atomicToken)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: resultCode,
		Message:    message,
		Data: responser.ShortUrlClickInfos{
			ShortUrlClickInfo: data,
			ClickAmount:       len(data),
		},
		TimeStamp: time.Now(),
	})
}
