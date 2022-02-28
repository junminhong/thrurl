package repo

import (
	"github.com/junminhong/thrurl/domain"
	"gorm.io/gorm"
)

type urlRepo struct {
	db *gorm.DB
}

func (urlRepo urlRepo) GetShortUrl(trackerID string) (shortUrl domain.ShortUrl, err error) {
	err = urlRepo.db.Where("tracker_id = ?", trackerID).First(&shortUrl).Error
	return shortUrl, err
}

func (urlRepo *urlRepo) SaveShortUrlClickInfo(shortUrl domain.ShortUrl, clickInfos []domain.ClickInfo) error {
	shortUrl.ClickInfo = clickInfos
	return urlRepo.db.Save(&shortUrl).Error
}

func NewUrlRepo(db *gorm.DB) domain.UrlRepository {
	return &urlRepo{db: db}
}
