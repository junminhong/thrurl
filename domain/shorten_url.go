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
	MemberID       int
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

type ShortUrlResponse struct {
	ShortenID      string    `json:"shorten_id"`
	Source         string    `json:"source"`
	SourceB        string    `json:"source_b"`
	SourceBPercent string    `json:"source_b_percent"`
	WhoClick       bool      `json:"who_click"`
	ClickCount     int64     `json:"click_count"`
	Expired        time.Time `json:"expired"`
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
	StoreShortenUrl() error
}

type ShortenUrlUseCase interface {
	// ShortenUrl 建立短網址
	ShortenUrl(request requester.ShortenUrl) responser.Response
}
