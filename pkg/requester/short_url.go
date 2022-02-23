package requester

type ShortenUrl struct {
	SourceUrl string `json:"source_url" binding:"required"`
}

type EditShortUrl struct {
	TrackerID  string `json:"tracker_id" binding:"required"`
	SourceUrlA string `json:"source_url_a" binding:"required"`
	SourceUrlB string `json:"source_url_b"`
	ABPercent  int    `json:"ab_percent"`
	WhoClick   bool   `json:"who_click"`
}

type RecordShortUrlInfo struct {
	ClickerIP      string `json:"clicker_ip"`
	ClickerCountry string `json:"clicker_country"`
	ClickerCity    string `json:"clicker_city"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	Platform       string `json:"platform"`
	OS             string `json:"os"`
}
