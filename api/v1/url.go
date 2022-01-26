package v1

import (
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/junminhong/thrurl/db/blot"
	"github.com/junminhong/thrurl/db/mongo"
	"github.com/junminhong/thrurl/db/postgresql"
	"github.com/junminhong/thrurl/db/redis"
	"github.com/junminhong/thrurl/grpc"
	"github.com/junminhong/thrurl/grpc/proto"
	"github.com/junminhong/thrurl/model"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/mssola/user_agent"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var blotDB = blot.Setup()
var mongoDB = mongo.Setup()
var redisDB = redis.Setup()
var postgresDB = postgresql.Setup()

const (
	IPNotLoginLimit = 100
	IPLoginLimit    = 500
	ShortenIDLen    = 6
)

type ShortUrlReqByNotLogin struct {
	SourceUrl string `json:"source_url" binding:"required"`
	Expired   string `json:"expired"`
}

type ShortUrlReqByLogin struct {
	SourceUrl   string `json:"source_url" binding:"required"`
	Expired     string `json:"expired"`
	SourceUrlB  string `json:"source_url_b"`
	BUrlPercent string `json:"b_url_percent"`
	WhoClick    bool   `json:"who_click"`
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

type shortUrlReq struct {
	SourceUrl   string    `json:"source_url" binding:"required"`
	SourceUrlB  string    `json:"source_url_b"`
	BUrlPercent string    `json:"b_url_percent"`
	WhoClick    bool      `json:"who_click"`
	Expired     time.Time `json:"expired"`
}

func ShortUrl(c *gin.Context) {
	request := &shortUrlReq{}
	err := c.BindJSON(request)
	if err != nil {
		// request bind data error
		log.Println(err.Error())
		log.Println(request)
		c.JSON(handler.BadRequest, handler.Response{
			ResultCode: handler.BadRequest,
			Message:    handler.ResponseFlag[handler.BadRequest],
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	isLife := handler.UrlLifeCheck(request.SourceUrl)
	if !isLife {
		c.JSON(handler.BadRequest, handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "無效的網址連結",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	callTimes, err := redisDB.Get(context.Background(), c.ClientIP()).Int()
	redisDB.Set(context.Background(), c.ClientIP(), callTimes+1, 60*time.Second)
	// bind request data後先判斷一下這筆需求是不是有帶token先區分登入狀態
	auth := c.Request.Header.Get("Authorization")
	if auth != "" {
		// 此筆請求有帶Authorization
		// 做字串切割，把token分出來
		tokens := strings.Split(auth, "Bearer ")
		if len(tokens) == 2 {
			// 判斷切割後的token資料是正確，正確切割長度應該為2
			token := tokens[1]
			shortenID := loginHandler(token, request)
			if shortenID == "" {
				c.JSON(handler.BadRequest, handler.Response{
					ResultCode: handler.BadRequest,
					Message:    "token失效",
					Data:       "",
					TimeStamp:  time.Now().UTC(),
				})
				return
			}
			c.JSON(handler.OK, handler.Response{
				ResultCode: handler.OK,
				Message:    "短網址生成完成",
				Data:       data{ShortUrl: os.Getenv("HOST_NAME") + "/" + shortenID},
				TimeStamp:  time.Now().UTC(),
			})
			return
		}
	}
	// 沒有登入
	if callTimes >= IPNotLoginLimit {
		c.JSON(handler.BadRequest, handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "你已經達到今天的短網址製作上限嘍",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	shortenID := notLoginHandler(request.SourceUrl)
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.OK,
		Message:    "短網址生成完成",
		Data:       data{ShortUrl: os.Getenv("HOST_NAME") + "/" + shortenID},
		TimeStamp:  time.Now().UTC(),
	})
	return
}

func getMemberIDByGrpc(token string) string {
	conn := grpc.SetupClient()
	conn.GetMethodConfig("New")

	client := proto.NewTokenAuthServiceClient(conn)
	result, err := client.VerifyAccessToken(context.Background(), &proto.TokenAuthRequest{Token: token})
	if err != nil {
		log.Println(err.Error())
	}
	if result != nil {
		return result.MemberID
	}
	return ""
}

func loginHandler(token string, request *shortUrlReq) string {
	tmp := getMemberIDByGrpc(token)
	if tmp == "" {
		return ""
	}
	memberID, err := strconv.Atoi(tmp)
	malice, maliceType := handler.SafeUrlCheck(request.SourceUrl)
	var (
		maliceB     bool
		maliceTypeB string
	)
	if request.SourceUrlB != "" {
		maliceB, maliceTypeB = handler.SafeUrlCheck(request.SourceUrl)
	}
	data := model.ShortUrl{
		MemberID:       memberID,
		Source:         request.SourceUrl,
		SourceB:        request.SourceUrlB,
		SourceBPercent: request.BUrlPercent,
		Malice:         malice,
		MaliceType:     maliceType,
		MaliceB:        maliceB,
		MaliceTypeB:    maliceTypeB,
		WhoClick:       request.WhoClick,
		Expired:        request.Expired,
	}
	postgresDB.Create(&data)
	data.ShortenID = shortenIDHandler(data.ID)
	err = postgresDB.Save(&data).Error
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"short_urls_shorten_id_key\" (SQLSTATE 23505)" {
			data.ShortenID = shortenIDHandler(data.ID)
			postgresDB.Save(&data)
		}
	}
	insertBlot(data.ShortenID+",member", data.Source)
	return data.ShortenID
}

func shortenIDHandler(id int) string {
	base62 := handler.Encode(id)
	salt := handler.GetSalt(ShortenIDLen - len(base62))
	return salt + base62
}

func notLoginHandler(sourceUrl string) string {
	//沒有登入的狀態處理比較簡單，只需要做到縮短網址就好了
	result, maliceType := handler.SafeUrlCheck(sourceUrl)
	data := model.ShortUrl{Source: sourceUrl, Malice: result, MaliceType: maliceType}
	postgresDB.Create(&data)
	data.ShortenID = shortenIDHandler(data.ID)
	err := postgresDB.Save(&data).Error
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"short_urls_shorten_id_key\" (SQLSTATE 23505)" {
			data.ShortenID = shortenIDHandler(data.ID)
			postgresDB.Save(&data)
		}
	}
	insertBlot(data.ShortenID, data.Source)
	return data.ShortenID
}

func insertBlot(shortenID string, sourceUrl string) {
	err := blotDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UrlBucket"))
		err := bucket.Put([]byte(shortenID), []byte(sourceUrl))
		return err
	})
	if err != nil {
		log.Println(err.Error())
	}
}

// @Summary 縮短網址
// @Tags url
// @version 1.0
// @Accept application/json
// @produce application/json
// @param data body ShortUrlReqByNotLogin true "縮短網址請求資料"
// @Success 200 {object} handler.Response
// @Router /short-url [post]
func ShortUrlByNotLogin(ctx *gin.Context) {
	request := &ShortUrlReqByNotLogin{}
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
	response, err = shorUrlNotLoginHandler(request)
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(response.ResultCode, response)
}

// @Summary 縮短網址(有會員)
// @Tags url
// @version 1.0
// @Accept application/json
// @produce application/json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @param data body ShortUrlReqByLogin true "縮短網址請求資料"
// @Success 200 {object} handler.Response
// @Router /short-url [post]
/*func shortUrlByLogin() {
	auth := ctx.Request.Header.Get("Authorization")
	if auth != "" {

	}
	// 代表有登入
	tokenAry := strings.Split(auth, "Bearer ")
	if len(tokenAry) != 2 {
		response = handler.Response{
			ResultCode: handler.Forbidden,
			Message:    handler.ResponseFlag[handler.Forbidden],
			Data:       &data{ShortUrl: ""},
			TimeStamp:  time.Now().UTC(),
		}
		ctx.JSON(response.ResultCode, response)
		return
	}
	rpc(tokenAry[1], request)
}*/

func rpc(token string, request *ShortUrlReqByLogin) {
	conn := grpc.SetupClient()
	conn.GetMethodConfig("New")

	client := proto.NewTokenAuthServiceClient(conn)
	result, err := client.VerifyAccessToken(context.Background(), &proto.TokenAuthRequest{Token: token})
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(result.MemberID)
	collection := mongoDB.Database("thrurl").Collection("record")
	res, _ := collection.InsertOne(context.Background(), bson.M{
		"member_id":    result.MemberID,
		"source_a":     request.SourceUrl,
		"source_b":     request.SourceUrlB,
		"click_amount": 0,
	})
	id := res.InsertedID
	log.Println(id)
}

func shorUrlLoginHandler(token string, request *ShortUrlReqByLogin) (handler.Response, error) {
	// 處理過期、AB測試、點擊成效
	response := handler.Response{
		ResultCode: handler.OK,
		Message:    handler.ResponseFlag[handler.OK],
		Data:       data{ShortUrl: ""},
		TimeStamp:  time.Now().UTC(),
	}
	return response, nil
}

func shorUrlNotLoginHandler(request *ShortUrlReqByNotLogin) (handler.Response, error) {
	/*for index := 1; index < 1000; index++ {
		base62 := Encode(index)
		salt := getTimeNano(6 - len(base62))
		log.Println(base62 + salt)
	}*/
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
	// 用blot比較快
	sourceUrl := ""
	_ = blotDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UrlBucket"))
		sourceUrl = string(bucket.Get([]byte(ctx.Param("short-id"))))
		if sourceUrl == "" {
			// 如果是會員的要去抓db會有點時間吧
			sourceUrl = string(bucket.Get([]byte(ctx.Param("short-id") + ",member")))
			shortUrl := model.ShortUrl{}
			postgresDB.Where("shorten_id = ?", ctx.Param("short-id")).First(&shortUrl)
			if shortUrl.WhoClick {
				ua := user_agent.New(ctx.GetHeader("User-Agent"))
				browserName, browserVersion := ua.Browser()
				shortUrlInfo := model.ShortUrlInfo{
					ClickerIP:      ctx.ClientIP(),
					Browser:        browserName,
					BrowserVersion: browserVersion,
					Platform:       ua.Platform(),
					OS:             ua.OS(),
				}
				shortUrl.ShortUrlInfos = []model.ShortUrlInfo{shortUrlInfo}
				postgresDB.Save(&shortUrl)
			}
			if shortUrl.SourceB != "" {
				// 有ab測試
				log.Println("ab測試")
				percent, _ := strconv.Atoi(shortUrl.SourceBPercent)
				sourceUrl = handler.ABTest(shortUrl.Source, shortUrl.SourceB, percent)
				log.Println(sourceUrl)
			}
		}
		return nil
	})
	ctx.Redirect(http.StatusFound, sourceUrl)
}
func shortInfoHandler() {

}

func TestMongo(ctx *gin.Context) {
	collection := mongoDB.Database("thrurl").Collection("record")
	/*cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  res, _ := collection.InsertOne(cntx, bson.M{"name": "pi", "value": 3.14159})
	  id := res.InsertedID
	  log.Println(id)*/

	var result struct {
		Name  string
		Value float64
	}

	//objID, _ := primitive.ObjectIDFromHex("61e9807cf8d9a401fbf98275")
	filter := bson.M{"name": "pi"}
	cntx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(cntx, filter).Decode(&result)
	if err != nil {
		// Do something when no record was found
		fmt.Println("record does not exist")
	}
	log.Println(result)
}
