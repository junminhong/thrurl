package application

import (
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/base62"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"strings"
)

type shortenUrlApp struct {
	shortenUrlRepo domain.ShortUrlRepository
}

func NewShortenUrlUseCase(shortenUrlRepo domain.ShortUrlRepository) domain.ShortUrlApp {
	return &shortenUrlApp{shortenUrlRepo}
}

func (shortenUrlApp *shortenUrlApp) ShortenUrl(source string, atomicToken string) (resultCode int, message string, data responser.ShortenUrl) {
	if urlLife := handler.UrlLifeCheck(source); !urlLife {
		return responser.UrlLinkNotFoundErr.Code(), responser.UrlLinkNotFoundErr.Message(), responser.ShortenUrl{}
	}
	// 跟repo說要存東西吧
	shortUrlId, _ := shortenUrlApp.shortenUrlRepo.GetShortUrlLastID()
	shortUrlInfo := domain.ShortUrlInfo{SourceUrlA: source}
	trackerID := base62.GetSaltEncode(shortUrlId, 6)
	shortUrl := domain.ShortUrl{ShortUrlInfo: shortUrlInfo, TrackerID: trackerID}
	if atomicToken != "" {
		memberUUID, _ := shortenUrlApp.shortenUrlRepo.GetMemberUUID(atomicToken)
		shortUrl.MemberID = memberUUID
	}
	if err := shortenUrlApp.shortenUrlRepo.SaveShortUrl(shortUrl); err != nil {
		return responser.SaveShortUrlErr.Code(), responser.SaveShortUrlErr.Message(), responser.ShortenUrl{}
	}
	return responser.SaveShortUrlOk.Code(), responser.SaveShortUrlOk.Message(), responser.ShortenUrl{
		ShortenUrl: "http://127.0.0.1:9220/" + trackerID,
	}
}

func (shortenUrlApp *shortenUrlApp) GetSourceUrl(trackerId string) (sourceUrl string) {
	sourceUrl, _ = shortenUrlApp.shortenUrlRepo.GetSourceUrl(trackerId)
	return sourceUrl
}

func (shortenUrlApp *shortenUrlApp) EditShortUrl(editShortUrl requester.EditShortUrl, atomicToken string) (resultCode int, message string) {
	memberUUID, _ := shortenUrlApp.shortenUrlRepo.GetMemberUUID(atomicToken)
	if memberUUID == "" {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message()
	}
	shortUrl, err := shortenUrlApp.shortenUrlRepo.GetShortUrl(editShortUrl.TrackerID)
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message()
	}
	if strings.Compare(shortUrl.MemberID, memberUUID) != 0 {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message()
	}
	shortUrlInfo, err := shortenUrlApp.shortenUrlRepo.GetShortUrlInfo(shortUrl.ID)
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message()
	}
	shortUrl.WhoClick = editShortUrl.WhoClick
	shortUrlInfo.SourceUrlA = editShortUrl.SourceUrlA
	shortUrlInfo.SourceUrlB = editShortUrl.SourceUrlB
	shortUrlInfo.ABPercent = editShortUrl.ABPercent
	if err := shortenUrlApp.shortenUrlRepo.SaveShortUrl(shortUrl); err != nil {
		return responser.SaveShortUrlErr.Code(), responser.SaveShortUrlErr.Reload("short url更新失敗").Message()
	}
	if err := shortenUrlApp.shortenUrlRepo.SaveShortUrlInfo(shortUrlInfo); err != nil {
		return responser.SaveShortUrlErr.Code(), responser.SaveShortUrlErr.Reload("short url更新失敗").Message()
	}
	return responser.SaveShortUrlOk.Code(), responser.SaveShortUrlOk.Reload("short url更新成功").Message()
}
