package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/interfaces/grpc/proto"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

type shortenUrlHandler struct {
	grpcClient        *grpc.ClientConn
	shortenUrlUseCase domain.ShortenUrlUseCase
	shortenUrlRepo    domain.ShortenUrlRepository
}

func NewShortenUrlHandler(router *gin.Engine, grpcClient *grpc.ClientConn, shortenUrlUseCase domain.ShortenUrlUseCase, shortenUrlRepo domain.ShortenUrlRepository) {
	handler := &shortenUrlHandler{grpcClient, shortenUrlUseCase, shortenUrlRepo}
	router.POST("/api/v1/url", handler.ShortenUrl)
	router.GET("/:shorten-id", handler.RedirectShortenUrl)
	router.GET("/test", handler.test)
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
	isLife := handler.UrlLifeCheck(request.SourceUrl)
	if !isLife {
		c.JSON(http.StatusOK, responser.Response{
			ResultCode: responser.UrlLinkNotFoundErr.Code(),
			Message:    responser.UrlLinkNotFoundErr.Message(),
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	atomicToken := requester.GetAtomicToken(c)
	memberUUID := shortenUrlHandler.getMemberUUIDByGrpc(atomicToken)
	response := shortenUrlHandler.shortenUrlUseCase.ShortenUrl(request, memberUUID)
	c.JSON(http.StatusOK, response)
}

func (shortenUrlHandler *shortenUrlHandler) test(c *gin.Context) {
	atomicToken := requester.GetAtomicToken(c)
	log.Println(atomicToken)
	log.Println(shortenUrlHandler.getMemberUUIDByGrpc(atomicToken))
}

func (shortenUrlHandler *shortenUrlHandler) getMemberUUIDByGrpc(atomicToken string) string {
	if atomicToken != "" {
		client := proto.NewMemberServiceClient(shortenUrlHandler.grpcClient)
		result, err := client.VerifyAtomicToken(context.Background(), &proto.AtomicTokenAuthRequest{AtomicToken: atomicToken})
		if err != nil {
			log.Println(err.Error())
		}
		return result.MemberUUID
	}
	return ""
}

func (shortenUrlHandler *shortenUrlHandler) RedirectShortenUrl(c *gin.Context) {
	url := shortenUrlHandler.shortenUrlUseCase.GetShortenUrl(c.Param("shorten-id"))
	c.Redirect(http.StatusMovedPermanently, url)
}
