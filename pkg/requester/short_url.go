package requester

type ShortenUrl struct {
	SourceUrl string `json:"source_url" binding:"required"`
}

type EditShortUrl struct {
	TrackerID  string `json:"tracker_id" binding:"required"`
	SourceUrlA string `json:"source_url_a" binding:"required"`
	SourceUrlB string `json:"source_url_b"`
	ABPercent  string `json:"ab_percent"`
	WhoClick   bool   `json:"who_click"`
}
