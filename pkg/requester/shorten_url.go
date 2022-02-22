package requester

type ShortenUrl struct {
	SourceUrl string `json:"source_url" bind:"required"`
}
