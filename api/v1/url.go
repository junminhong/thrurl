package v1

import (
	"context"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/junminhong/thrurl/db/blot"
	"github.com/junminhong/thrurl/db/postgresql"
	"github.com/junminhong/thrurl/db/redis"
	"github.com/junminhong/thrurl/grpc"
	"github.com/junminhong/thrurl/grpc/proto"
	"github.com/junminhong/thrurl/model"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/mssola/user_agent"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var blotDB = blot.Setup()
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

// AllUrlList 取得某的會員的所有url資料
// 分頁問題？？
func AllUrlList(c *gin.Context) {
	token := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]
	memberID := getMemberIDByGrpc(token)
	if memberID == "" {
		// 沒有辦法取得member id 代表該token不是合法的或已經失效了
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.TokenError2,
			Message:    handler.ErrorFlag[handler.TokenError2],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	// 有取得memberID
	rows, err := postgresDB.Model(&model.ShortUrl{}).Where("member_id = ?", memberID).Rows()
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	shortUrls := make([]model.ShortUrlResponse, 0)
	for rows.Next() {
		shortUrl := model.ShortUrlResponse{}
		postgresDB.ScanRows(rows, &shortUrl)
		shortUrls = append(shortUrls, shortUrl)
	}
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.OK,
		Message:    handler.ResponseFlag[handler.OK],
		Data:       shortUrls,
		TimeStamp:  time.Now().UTC(),
	})
}
