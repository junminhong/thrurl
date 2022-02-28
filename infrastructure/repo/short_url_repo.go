package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/junminhong/thrurl/domain"
	"github.com/junminhong/thrurl/infrastructure/grpc/proto"
	"github.com/junminhong/thrurl/pkg/responser"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type shortUrlRepo struct {
	db         *gorm.DB
	redis      *redis.Client
	grpcClient *grpc.ClientConn
}

func (shortUrlRepo shortUrlRepo) GetShortUrlListCount(memberUUID string) (count int64, err error) {
	shortUrl := domain.ShortUrl{}
	err = shortUrlRepo.db.Model(&shortUrl).Where("member_id=?", memberUUID).Count(&count).Error
	return count, err
}

func (shortUrlRepo *shortUrlRepo) GetShortUrlClickInfo(shortUrl domain.ShortUrl) (shortUrlClickInfos []responser.ShortUrlClickInfo, err error) {
	clickInfo := domain.ClickInfo{}
	rows, err := shortUrlRepo.db.Model(&clickInfo).Where("short_url_id = ?", shortUrl.ID).Rows()
	shortUrlClickInfo := responser.ShortUrlClickInfo{}
	for rows.Next() {
		shortUrlRepo.db.ScanRows(rows, &shortUrlClickInfo)
		shortUrlClickInfos = append(shortUrlClickInfos, shortUrlClickInfo)
	}
	return shortUrlClickInfos, err
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
	if err = shortUrlRepo.db.Model(&shortUrl).Association("ShortUrlInfo").Find(&shortUrlInfo); err != nil {
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

func (shortUrlRepo shortUrlRepo) GetShortUrlList(memberUUID string, limit int, offset int) (shortUrlLists []responser.ShortUrlList, err error) {
	shortUrl := domain.ShortUrl{}
	// rows, err := shortUrlRepo.db.Model(&shortUrl).Joins("ShortUrlInfo").Where("member_id = ?", memberUUID).Limit(limit).Offset(offset).Rows()
	rows, err := shortUrlRepo.db.Model(&shortUrl).Order("short_urls.id asc").Select("*").Joins("INNER JOIN short_url_infos on short_url_infos.short_url_id=short_urls.id").Where("member_id=?", memberUUID).Limit(limit).Offset(offset).Rows()
	defer rows.Close()
	shortUrlList := responser.ShortUrlList{}
	for rows.Next() {
		shortUrlRepo.db.ScanRows(rows, &shortUrlList)
		shortUrlLists = append(shortUrlLists, shortUrlList)
	}
	return shortUrlLists, err
}

func (shortUrlRepo shortUrlRepo) QuickGetShortUrlInfo(trackerID string) (shortUrlInfo domain.ShortUrlInfo, err error) {
	shortUrl := domain.ShortUrl{}
	if err = shortUrlRepo.db.Where("tracker_id=?", trackerID).First(&shortUrl).Error; err != nil {
		return shortUrlInfo, err
	}
	if err = shortUrlRepo.db.Model(&shortUrl).Association("ShortUrlInfo").Find(&shortUrlInfo); err != nil {
		return shortUrlInfo, err
	}
	return shortUrlInfo, err
}
