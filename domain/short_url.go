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
	ABPercent      string
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
	GetShortUrlInfo(shortUrlID int) (shortUrlInfo ShortUrlInfo, err error)
}

type ShortUrlApp interface {
	ShortenUrl(source string, atomicToken string) (resultCode int, message string, data responser.ShortenUrl)
	GetSourceUrl(trackerId string) (sourceUrl string)
	EditShortUrl(editShortUrl requester.EditShortUrl, atomicToken string) (resultCode int, message string)
}
