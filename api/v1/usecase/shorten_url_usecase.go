package usecase

import (
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/handler"
	"github.com/junminhong/thrurl/pkg/requester"
	"github.com/junminhong/thrurl/pkg/responser"
	"time"
)

type shortenUrlUseCase struct {
	shortenUrlRepo domain.ShortenUrlRepository
}

func NewShortenUrlUseCase(shortenUrlRepo domain.ShortenUrlRepository) domain.ShortenUrlUseCase {
	return &shortenUrlUseCase{shortenUrlRepo}
}

func (shortenUrlUseCase *shortenUrlUseCase) ShortenUrl(request requester.ShortenUrl) responser.Response {
	shortenUrl, err := shortenUrlUseCase.shortenUrlRepo.CreateShortenUrl()
	if err != nil {
		return responser.Response{
			ResultCode: responser.StoreShortenUrlErr.Code(),
			Message:    responser.StoreShortenUrlErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		}
	}
	base62 := handler.Encode(shortenUrl.ID)
	saltBase62 := handler.GetSalt(6 - len(base62))
	shortenUrl.Source = request.SourceUrl
	shortenUrl.ShortenID = saltBase62 + base62
	if err := shortenUrlUseCase.shortenUrlRepo.UpdateShortenUrl(shortenUrl); err != nil {
		return responser.Response{
			ResultCode: responser.StoreShortenUrlErr.Code(),
			Message:    responser.StoreShortenUrlErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		}
	}
	return responser.Response{
		ResultCode: responser.StoreShortenUrlOk.Code(),
		Message:    responser.StoreShortenUrlOk.Message(),
		Data:       responser.ShortenUrl{ShortenUrl: "http://127.0.0.1:9220/" + shortenUrl.ShortenID},
		TimeStamp:  time.Now(),
	}
}

func (shortenUrlUseCase *shortenUrlUseCase) GetShortenUrl(shortenID string) string {
	shortenUrl := shortenUrlUseCase.shortenUrlRepo.GetUrlByShortenID(shortenID)
	return shortenUrl.Source
}
