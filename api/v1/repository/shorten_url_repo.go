package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/pkg/handler"
	"gorm.io/gorm"
)

type shortenUrlRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewShortenUrlRepo(db *gorm.DB, redis *redis.Client) domain.ShortenUrlRepository {
	return &shortenUrlRepo{db, redis}
}

func (shortenUrlRepo *shortenUrlRepo) StoreShortenUrl() error {
	shortenUrl := domain.ShortenUrl{}
	if err := shortenUrlRepo.db.Create(&shortenUrl).Error; err != nil {
		return shortenUrlRepo.db.Create(&shortenUrl).Error
	}
	base62 := handler.Encode(shortenUrl.ID)
	saltBase62 := handler.GetSalt(len(base62))
	shortenUrl.ShortenID = saltBase62 + base62
	return shortenUrlRepo.db.Save(&shortenUrl).Error
}
