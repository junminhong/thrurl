package model

import (
	"gorm.io/gorm"
	"time"
)

type ShortUrl struct {
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
	ShortUrlInfos  []ShortUrlInfo `gorm:"foreignKey:ShortUrlID"`
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

type ShortUrlInfo struct {
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
