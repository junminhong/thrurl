package application

import (
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/base62"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"github.com/spf13/viper"
	"strings"
)

type shortUrlApp struct {
	shortUrlRepo domain.ShortUrlRepository
}

func (shortUrlApp *shortUrlApp) GetShortUrlClickInfo(trackerID string, atomicToken string) (resultCode int, message string, data []responser.ShortUrlClickInfo) {
	memberUUID, _ := shortUrlApp.shortUrlRepo.GetMemberUUID(atomicToken)
	if memberUUID == "" {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message(), []responser.ShortUrlClickInfo{}
	}
	shortUrl, err := shortUrlApp.shortUrlRepo.GetShortUrl(trackerID)
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message(), []responser.ShortUrlClickInfo{}
	}
	if strings.Compare(shortUrl.MemberID, memberUUID) != 0 {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message(), []responser.ShortUrlClickInfo{}
	}
	clickInfos, _ := shortUrlApp.shortUrlRepo.GetShortUrlClickInfo(shortUrl)
	return 0, "", clickInfos
}

func NewShortenUrlUseCase(shortenUrlRepo domain.ShortUrlRepository) domain.ShortUrlApp {
	return &shortUrlApp{shortenUrlRepo}
}

func (shortUrlApp *shortUrlApp) ShortenUrl(source string, atomicToken string) (resultCode int, message string, data responser.ShortenUrl) {
	if urlLife := handler.UrlLifeCheck(source); !urlLife {
		return responser.UrlLinkNotFoundErr.Code(), responser.UrlLinkNotFoundErr.Message(), responser.ShortenUrl{}
	}
	// 跟repo說要存東西吧
	shortUrlId, _ := shortUrlApp.shortUrlRepo.GetShortUrlLastID()
	shortUrlInfo := domain.ShortUrlInfo{SourceUrlA: source}
	trackerID := base62.GetSaltEncode(shortUrlId, 6)
	shortUrl := domain.ShortUrl{ShortUrlInfo: shortUrlInfo, TrackerID: trackerID}
	if atomicToken != "" {
		memberUUID, _ := shortUrlApp.shortUrlRepo.GetMemberUUID(atomicToken)
		shortUrl.MemberID = memberUUID
	}
	if err := shortUrlApp.shortUrlRepo.SaveShortUrl(shortUrl); err != nil {
		return responser.SaveShortUrlErr.Code(), responser.SaveShortUrlErr.Message(), responser.ShortenUrl{}
	}
	return responser.SaveShortUrlOk.Code(), responser.SaveShortUrlOk.Message(), responser.ShortenUrl{
		ShortenUrl: viper.GetString("APP.SHORT_URL_HOST") + "/" + trackerID,
	}
}

func (shortUrlApp *shortUrlApp) GetSourceUrl(trackerID string) (resultCode int, message string, sourceUrl string) {
	shortUrlInfo, _ := shortUrlApp.shortUrlRepo.QuickGetShortUrlInfo(trackerID)
	if shortUrlInfo.ABPercent == 0 {
		return responser.GetSourceUrlOk.Code(), responser.GetSourceUrlOk.Message(), shortUrlInfo.SourceUrlA
	}
	return responser.GetSourceUrlOk.Code(), responser.GetSourceUrlOk.Message(), handler.ABTest(shortUrlInfo.SourceUrlA, shortUrlInfo.SourceUrlB, shortUrlInfo.ABPercent)
}

func (shortUrlApp *shortUrlApp) EditShortUrl(editShortUrl requester.EditShortUrl, atomicToken string) (resultCode int, message string) {
	memberUUID, _ := shortUrlApp.shortUrlRepo.GetMemberUUID(atomicToken)
	if memberUUID == "" {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message()
	}
	shortUrl, err := shortUrlApp.shortUrlRepo.GetShortUrl(editShortUrl.TrackerID)
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message()
	}
	if strings.Compare(shortUrl.MemberID, memberUUID) != 0 {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message()
	}
	shortUrlInfo, err := shortUrlApp.shortUrlRepo.GetShortUrlInfo(shortUrl.ID)
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message()
	}
	shortUrl.WhoClick = editShortUrl.WhoClick
	shortUrlInfo.SourceUrlA = editShortUrl.SourceUrlA
	shortUrlInfo.SourceUrlB = editShortUrl.SourceUrlB
	shortUrlInfo.ABPercent = editShortUrl.ABPercent
	if err := shortUrlApp.shortUrlRepo.SaveShortUrl(shortUrl); err != nil {
		return responser.SaveShortUrlErr.Code(), responser.SaveShortUrlErr.Reload("short url更新失敗").Message()
	}
	if err := shortUrlApp.shortUrlRepo.SaveShortUrlInfo(shortUrlInfo); err != nil {
		return responser.SaveShortUrlErr.Code(), responser.SaveShortUrlErr.Reload("short url更新失敗").Message()
	}
	return responser.SaveShortUrlOk.Code(), responser.SaveShortUrlOk.Reload("short url更新成功").Message()
}

func (shortUrlApp *shortUrlApp) GetShortUrlList(limit int, offset int, atomicToken string) (resultCode int, message string, data []responser.ShortUrlList, page int) {
	memberUUID, _ := shortUrlApp.shortUrlRepo.GetMemberUUID(atomicToken)
	if memberUUID == "" {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message(), []responser.ShortUrlList{}, 0
	}
	shortUrlLists, _ := shortUrlApp.shortUrlRepo.GetShortUrlList(memberUUID, limit, offset)
	page = len(shortUrlLists) / limit
	if len(shortUrlLists)%limit != 0 {
		page += 1
	}
	return responser.GetShortUrlListOk.Code(), responser.GetShortUrlListOk.Message(), shortUrlLists, page
}

func (shortUrlApp *shortUrlApp) GetShortUrl(trackerID string, atomicToken string) (resultCode int, message string, data responser.ShortUrlInfo) {
	shortUrl, err := shortUrlApp.shortUrlRepo.GetShortUrl(trackerID)
	memberUUID, _ := shortUrlApp.shortUrlRepo.GetMemberUUID(atomicToken)
	if memberUUID == "" || strings.Compare(shortUrl.MemberID, memberUUID) != 0 {
		return responser.NotFoundAtomicTokenErr.Code(), responser.NotFoundAtomicTokenErr.Message(), responser.ShortUrlInfo{}
	}
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message(), responser.ShortUrlInfo{}
	}
	shortUrlInfo, _ := shortUrlApp.shortUrlRepo.QuickGetShortUrlInfo(trackerID)
	if err != nil {
		return responser.NotFoundShortUrlErr.Code(), responser.NotFoundShortUrlErr.Message(), responser.ShortUrlInfo{}
	}
	return responser.GetShortUrlListOk.Code(), responser.GetShortUrlListOk.Message(), responser.ShortUrlInfo{
		SourceUrlA: shortUrlInfo.SourceUrlA,
		SourceUrlB: shortUrlInfo.SourceUrlB,
		ABPercent:  shortUrlInfo.ABPercent,
		Expired:    shortUrl.Expired,
	}
}
