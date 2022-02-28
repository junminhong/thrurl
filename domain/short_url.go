package domain

import (
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"gorm.io/gorm"
	"time"
)

type ShortUrl struct {
	gorm.Model
	ID           int    `gorm:"primary_key;auto_increment;not_null"`
	TrackerID    string `gorm:"unique"`
	MemberID     string
	WhoClick     bool
	ShortUrlInfo ShortUrlInfo
	ClickInfo    []ClickInfo
	Expired      time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type ShortUrlInfo struct {
	gorm.Model
	ShortUrlID     int
	SourceUrlA     string
	SourceUrlB     string
	ABPercent      int
	IsMaliceUrlA   bool
	MaliceUrlAType string
	IsMaliceUrlB   bool
	MaliceUrlBType string
}

type ClickInfo struct {
	gorm.Model
	ShortUrlID     int
	ClickerIP      string
	ClickerCountry string
	ClickerCity    string
	Browser        string
	BrowserVersion string
	Platform       string
	OS             string
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type ShortUrlRepository interface {
	SaveShortUrl(shortUrl ShortUrl) error
	SaveShortUrlInfo(shortUrlInfo ShortUrlInfo) error
	GetMemberUUID(atomicToken string) (memberUUID string, err error)
	GetShortUrlLastID() (ID int, err error)
	GetSourceUrl(trackerID string) (sourceUrl string, err error)
	GetShortUrl(trackerID string) (shortUrl ShortUrl, err error)
	GetShortUrlList(memberUUID string, limit int, offset int) (shortUrlLists []responser.ShortUrlList, err error)
	GetShortUrlListCount(memberUUID string) (count int64, err error)
	GetShortUrlInfo(shortUrlID int) (shortUrlInfo ShortUrlInfo, err error)
	GetShortUrlClickInfo(shortUrl ShortUrl) (clickInfos []responser.ShortUrlClickInfo, err error)
	QuickGetShortUrlInfo(trackerID string) (shortUrlInfo ShortUrlInfo, err error)
}

type ShortUrlApp interface {
	ShortenUrl(source string, atomicToken string) (resultCode int, message string, data responser.ShortenUrl)
	GetSourceUrl(trackerID string) (resultCode int, message string, sourceUrl string)
	GetShortUrl(trackerID string, atomicToken string) (resultCode int, message string, data responser.ShortUrlInfo)
	GetShortUrlList(limit int, offset int, atomicToken string) (resultCode int, message string, data []responser.ShortUrlList, page int)
	GetShortUrlClickInfo(trackerID string, atomicToken string) (resultCode int, message string, data []responser.ShortUrlClickInfo)
	EditShortUrl(editShortUrl requester.EditShortUrl, atomicToken string) (resultCode int, message string)
}
