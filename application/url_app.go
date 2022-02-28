package application

import (
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
)

type urlApp struct {
	urlRepo domain.UrlRepository
}

func (urlApp *urlApp) RecordShortUrlClickInfo(trackerID string, request requester.RecordShortUrlInfo) (resultCode int, message string) {
	shortUrl, _ := urlApp.urlRepo.GetShortUrl(trackerID)
	if !shortUrl.WhoClick {
		return responser.RecordClickInfoErr.Code(), responser.RecordClickInfoErr.Reload("沒有開啟點擊成效，不需要記錄").Message()
	}
	clickInfo := domain.ClickInfo{
		ClickerIP:      request.ClickerIP,
		Browser:        request.Browser,
		BrowserVersion: request.BrowserVersion,
		OS:             request.OS,
		Platform:       request.Platform,
	}
	if err := urlApp.urlRepo.SaveShortUrlClickInfo(shortUrl, []domain.ClickInfo{clickInfo}); err != nil {
		return responser.RecordClickInfoErr.Code(), responser.RecordClickInfoErr.Message()
	}
	return responser.RecordClickInfoOk.Code(), responser.RecordClickInfoOk.Message()
}

func NewUrlApp(urlRepo domain.UrlRepository) domain.UrlApp {
	return &urlApp{urlRepo: urlRepo}
}
