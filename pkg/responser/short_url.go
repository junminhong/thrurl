package responser

import "time"

type ShortenUrl struct {
	ShortenUrl string `json:"shorten_url"`
}

type ShortenUrlInfo struct {
	ShortenID      string    `json:"shorten_id"`
	Source         string    `json:"source"`
	SourceB        string    `json:"source_b"`
	SourceBPercent string    `json:"source_b_percent"`
	WhoClick       bool      `json:"who_click"`
	ClickCount     int64     `json:"click_count"`
	Expired        time.Time `json:"expired"`
}
