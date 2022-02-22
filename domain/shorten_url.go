package domain

import (
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"gorm.io/gorm"
	"time"
)

type ShortenUrl struct {
	gorm.Model
	ID             int    `gorm:"primary_key;auto_increment;not_null"`
	ShortenID      string `gorm:"unique"`
	MemberID       string
	Source         string
	SourceB        string
	SourceBPercent string
	Malice         bool
	MaliceType     string
	MaliceB        bool
	MaliceTypeB    string
	WhoClick       bool
	ShortUrlInfos  []ShortenUrlInfo `gorm:"foreignKey:ShortUrlID"`
	Expired        time.Time
}

type ShortenUrlInfo struct {
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
}

type ShortenUrlRepository interface {
	CreateShortenUrl() (*ShortenUrl, error)
	UpdateShortenUrl(shortenUrl *ShortenUrl) error
	GetUrlByShortenID(shortenID string) (shortenUrl *ShortenUrl)
}

type ShortenUrlUseCase interface {
	// ShortenUrl 建立短網址
	ShortenUrl(request requester.ShortenUrl, memberUUID string) responser.Response
	GetShortenUrl(shortenID string) string
}
