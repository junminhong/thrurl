package responser

import "time"

type ShortenUrl struct {
	ShortenUrl string `json:"shorten_url"`
}

type ShortUrlInfo struct {
	SourceUrlA string    `json:"source_url_a"`
	SourceUrlB string    `json:"source_url_b"`
	ABPercent  int       `json:"ab_percent"`
	Expired    time.Time `json:"expired"`
}

type ShortUrlClickInfo struct {
	ClickerIP      string    `json:"clicker_ip"`
	ClickerCountry string    `json:"clicker_country"`
	ClickerCity    string    `json:"clicker_city"`
	Browser        string    `json:"browser"`
	BrowserVersion string    `json:"browser_version"`
	Platform       string    `json:"platform"`
	OS             string    `json:"os"`
	CreatedAt      time.Time `json:"created_at"`
}

type ShortUrlClickInfos struct {
	ShortUrlClickInfo []ShortUrlClickInfo `json:"short_url_click_info"`
	ClickAmount       int                 `json:"click_amount"`
}

type ShortUrlList struct {
	TrackerID string `json:"tracker_id"`
	WhoClick  bool   `json:"who_click"`
}

type ShortUrlLists struct {
	ShortUrlList []ShortUrlList `json:"short_url_list"`
	Page         int            `json:"page"`
}

type GetSourceUrl struct {
	SourceUrl string `json:"source_url"`
}
