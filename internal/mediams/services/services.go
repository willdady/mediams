package services

import (
	"github.com/willdady/mediams/internal/mediams/models"
)

// MediaService interface
type MediaService interface {
	CreateMedia(media *models.Media) error
	UpdateMedia(media *models.Media) error
	DeleteMedia(media *models.Media) error
	GetMedia(mediaID uint64) (models.Media, error)
	GetMedias(cursor string, userID string, tag string) ([]models.Media, string, error)
	AlbumExists(albumID uint64) (bool, error)
	CreateAlbum(album *models.Album) error
	UpdateAlbum(album *models.Album) error
	DeleteAlbum(album *models.Album) error
	GetAlbum(albumID uint64) (models.Album, error)
	GetAlbums(cursor string, userID string) ([]models.Album, string, error)
	GetTags() ([]string, error)
}
