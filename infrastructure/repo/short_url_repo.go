package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/infrastructure/grpc/proto"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type shortUrlRepo struct {
	db         *gorm.DB
	redis      *redis.Client
	grpcClient *grpc.ClientConn
}

func NewShortenUrlRepo(db *gorm.DB, redis *redis.Client, grpcClient *grpc.ClientConn) domain.ShortUrlRepository {
	return &shortUrlRepo{db, redis, grpcClient}
}

func (shortUrlRepo *shortUrlRepo) SaveShortUrl(shortUrl domain.ShortUrl) error {
	return shortUrlRepo.db.Save(&shortUrl).Error
}

func (shortUrlRepo *shortUrlRepo) GetMemberUUID(atomicToken string) (memberUUID string, err error) {
	client := proto.NewMemberServiceClient(shortUrlRepo.grpcClient)
	result, err := client.VerifyAtomicToken(context.Background(), &proto.AtomicTokenAuthRequest{AtomicToken: atomicToken})
	if err != nil {
		return "", err
	}
	return result.MemberUUID, nil
}

func (shortUrlRepo *shortUrlRepo) GetShortUrlLastID() (ID int, err error) {
	shortUrl := domain.ShortUrl{}
	err = shortUrlRepo.db.Last(&shortUrl).Error
	return shortUrl.ID, err
}

func (shortUrlRepo *shortUrlRepo) GetSourceUrl(trackerID string) (sourceUrl string, err error) {
	shortUrlInfo := domain.ShortUrlInfo{}
	shortUrl := domain.ShortUrl{}
	if err = shortUrlRepo.db.Where("tracker_id=?", trackerID).First(&shortUrl).Error; err != nil {
		return "", err
	}
	if err = shortUrlRepo.db.Where("short_url_id=?", shortUrl.ID).First(&shortUrlInfo).Error; err != nil {
		return "", err
	}
	return shortUrlInfo.SourceUrlA, err
}

func (shortUrlRepo *shortUrlRepo) SaveShortUrlInfo(shortUrlInfo domain.ShortUrlInfo) error {
	return shortUrlRepo.db.Save(&shortUrlInfo).Error
}

func (shortUrlRepo *shortUrlRepo) GetShortUrl(trackerID string) (shortUrl domain.ShortUrl, err error) {
	err = shortUrlRepo.db.Where("tracker_id = ?", trackerID).First(&shortUrl).Error
	return shortUrl, err
}

func (shortUrlRepo *shortUrlRepo) GetShortUrlInfo(shortUrlID int) (shortUrlInfo domain.ShortUrlInfo, err error) {
	err = shortUrlRepo.db.Where("short_url_id=?", shortUrlID).First(&shortUrlInfo).Error
	return shortUrlInfo, err
}
