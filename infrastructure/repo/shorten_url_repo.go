package repo

import (
	"github.com/go-redis/redis/v8"
	"github.com/junminhong/thrurl/domain"
	"gorm.io/gorm"
)

type shortenUrlRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewShortenUrlRepo(db *gorm.DB, redis *redis.Client) domain.ShortenUrlRepository {
	return &shortenUrlRepo{db, redis}
}

func (shortenUrlRepo *shortenUrlRepo) CreateShortenUrl() (*domain.ShortenUrl, error) {
	shortenUrl := domain.ShortenUrl{}
	if err := shortenUrlRepo.db.Create(&shortenUrl).Error; err != nil {
		return &shortenUrl, shortenUrlRepo.db.Create(&shortenUrl).Error
	}
	return &shortenUrl, nil
}

func (shortenUrlRepo *shortenUrlRepo) UpdateShortenUrl(shortenUrl *domain.ShortenUrl) error {
	return shortenUrlRepo.db.Save(shortenUrl).Error
}

func (shortenUrlRepo *shortenUrlRepo) GetUrlByShortenID(shortenID string) (shortenUrl *domain.ShortenUrl) {
	shortenUrlRepo.db.Where("shorten_id=?", shortenID).First(&shortenUrl)
	return shortenUrl
}
