package postgres

import (
	"encoding/base64"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/willdady/mediams/internal/errors"
	"github.com/willdady/mediams/internal/mediams/models"
	"github.com/willdady/mediams/internal/utils"
)

type MediaService struct {
	DB *gorm.DB
}

func NewMediaService(db *gorm.DB) *MediaService {
	return &MediaService{DB: db}
}

func (service *MediaService) CreateMedia(media *models.Media) error {
	if err := service.DB.Create(media).Error; err != nil {
		return err
	}
	return nil
}

func (service *MediaService) UpdateMedia(media *models.Media) error {
	if err := service.DB.Save(media).Error; err != nil {
		return err
	}
	return nil
}

func (service *MediaService) DeleteMedia(media *models.Media) error {
	if media.ID == 0 {
		return &errors.DeleteIsMissingID{}
	}
	if err := service.DB.Delete(media).Error; err != nil {
		return err
	}
	return nil
}

func (service *MediaService) GetMedia(mediaID uint64) (models.Media, error) {
	media := models.Media{}
	if service.DB.Where("id = ?", mediaID).First(&media).RecordNotFound() {
		return media, &errors.NotFound{}
	}
	return media, nil
}

func (service *MediaService) GetMedias(cursor string, userID string, tag string) ([]models.Media, string, error) {
	medias := []models.Media{}
	query := service.DB.Order("id desc")
	if cursor != "" {
		id, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return medias, "", &errors.CursorDecodingError{}
		}
		query = query.Where("id <= ?", id)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if tag != "" {
		query = query.Where("? = ANY(tags)", tag)
	}
	// Note we over-fetch by 1 so we can check if there are more items
	nextCursor := ""
	limit := 101
	query = query.Limit(limit)
	if err := query.Find(&medias).Error; err != nil {
		return medias, nextCursor, nil
	}
	if len(medias) == limit {
		lastItem := medias[len(medias)-1]
		nextCursor = utils.UintToBase64(lastItem.ID)
		medias = medias[:len(medias)-1]
	}
	return medias, nextCursor, nil
}

func (service *MediaService) AlbumExists(albumID uint64) (bool, error) {
	result := struct {
		Exists bool
	}{}
	if err := service.DB.Raw("SELECT EXISTS(SELECT 1 FROM albums WHERE id=?) as exists", albumID).Scan(&result).Error; err != nil {
		return false, err
	}
	return result.Exists, nil
}

func (service *MediaService) CreateAlbum(album *models.Album) error {
	if err := service.DB.Create(album).Error; err != nil {
		return err
	}
	return nil
}

func (service *MediaService) UpdateAlbum(album *models.Album) error {
	if err := service.DB.Save(album).Error; err != nil {
		return err
	}
	return nil
}

func (service *MediaService) DeleteAlbum(album *models.Album) error {
	if album.ID == 0 {
		return &errors.DeleteIsMissingID{}
	}
	if err := service.DB.Delete(album).Error; err != nil {
		return err
	}
	return nil
}

func (service *MediaService) GetAlbum(albumID uint64) (models.Album, error) {
	album := models.Album{}
	if service.DB.Where("id = ?", albumID).First(&album).RecordNotFound() {
		return album, &errors.NotFound{}
	}
	return album, nil
}

func (service *MediaService) GetAlbums(cursor string, userID string) ([]models.Album, string, error) {
	albums := []models.Album{}
	query := service.DB.Order("id desc")
	if cursor != "" {
		id, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return albums, "", &errors.CursorDecodingError{}
		}
		query = query.Where("id <= ?", id)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	// Note we over-fetch by 1 so we can check if there are more items
	nextCursor := ""
	limit := 101
	query = query.Limit(limit)
	if err := query.Find(&albums).Error; err != nil {
		return albums, nextCursor, err
	}
	if len(albums) == limit {
		lastItem := albums[len(albums)-1]
		nextCursor = utils.UintToBase64(lastItem.ID)
		albums = albums[:len(albums)-1]
	}
	return albums, nextCursor, nil
}

func (service *MediaService) GetTags() ([]string, error) {
	tags := pq.StringArray{}
	service.DB.Raw("SELECT array_agg(DISTINCT flattags) FROM medias, unnest(tags) as flattags").Row().Scan(&tags)
	return tags, nil
}
