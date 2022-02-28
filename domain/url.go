package domain

import "github.com/junminhong/thrurl/pkg/requester"

type UrlRepository interface {
	SaveShortUrlClickInfo(shortUrl ShortUrl, clickInfos []ClickInfo) error
	GetShortUrl(trackerID string) (shortUrl ShortUrl, err error)
}

type UrlApp interface {
	RecordShortUrlClickInfo(trackerID string, request requester.RecordShortUrlInfo) (resultCode int, message string)
}
