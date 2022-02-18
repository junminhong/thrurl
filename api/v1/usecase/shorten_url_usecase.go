package usecase

import (
	"github.com/junminhong/thrurl/domain"
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
	if err := shortenUrlUseCase.shortenUrlRepo.StoreShortenUrl(); err != nil {
		// 儲存失敗
		return responser.Response{
			ResultCode: responser.StoreShortenUrlErr.Code(),
			Message:    responser.StoreShortenUrlErr.Message(),
			Data:       "",
			TimeStamp:  time.Now(),
		}
	}
	return responser.Response{
		ResultCode: responser.StoreShortenUrlErr.Code(),
		Message:    responser.StoreShortenUrlErr.Message(),
		Data:       "",
		TimeStamp:  time.Now(),
	}
}
