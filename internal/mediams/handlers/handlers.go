package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/willdady/mediams/internal/errors"
	"github.com/willdady/mediams/internal/mediams/models"
	"github.com/willdady/mediams/internal/mediams/services"
)

func handleServiceError(err error, c *gin.Context) {
	switch err.(type) {
	case *errors.NotFound:
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{"status": http.StatusNotFound, "message": err.Error()})
	case *errors.DeleteIsMissingID:
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
	case *errors.CursorDecodingError:
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
	default:
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"status": http.StatusInternalServerError, "message": "An internal server error occurred"})
	}
}

func getMediaServiceFromContext(c *gin.Context) services.MediaService {
	value, _ := c.Get("mediaService")
	return value.(services.MediaService)
}

func NotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Not found"})
}

func CreateMedia(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	media := &models.Media{}
	err := c.ShouldBindJSON(media)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	// If AlbumID is set, confirm album exists
	if media.AlbumID != 0 {
		albumExists, err := mediaService.AlbumExists(uint64(media.AlbumID))
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "message": err.Error()})
			return
		}
		if !albumExists {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "message": "Album matching id does not exist"})
			return
		}
	}
	err = mediaService.CreateMedia(media)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusCreated, media)
}

func UpdateMedia(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	mediaID := uint64(c.GetInt64("ID"))
	media := &models.Media{}
	err := c.ShouldBindJSON(media)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	var existingMedia models.Media
	existingMedia, err = mediaService.GetMedia(mediaID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	media.ID = uint(mediaID)
	media.CreatedAt = existingMedia.CreatedAt
	// If AlbumID is set, confirm album exists
	if media.AlbumID != 0 {
		albumExists, err := mediaService.AlbumExists(uint64(media.AlbumID))
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "message": err.Error()})
			return
		}
		if !albumExists {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"status": http.StatusBadRequest, "message": "Album matching id does not exist"})
			return
		}
	}
	err = mediaService.UpdateMedia(media)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, media)
}

func DeleteMedia(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	mediaID := uint64(c.GetInt64("ID"))
	media, err := mediaService.GetMedia(mediaID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	if err := mediaService.DeleteMedia(&media); err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetMedia(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	mediaID := uint64(c.GetInt64("ID"))
	media, err := mediaService.GetMedia(mediaID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, media)
}

func GetMedias(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	cursor := c.Query("cursor")
	userID := c.Query("userId")
	tag := c.Query("tag")
	medias, nextCursor, err := mediaService.GetMedias(cursor, userID, tag)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"nextCursor": nextCursor, "results": medias})
}

func CreateAlbum(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	album := &models.Album{}
	err := c.ShouldBindJSON(album)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	err = mediaService.CreateAlbum(album)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusCreated, album)
}

func UpdateAlbum(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	albumID := uint64(c.GetInt64("ID"))
	album := &models.Album{}
	err := c.ShouldBindJSON(album)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	var existingAlbum models.Album
	existingAlbum, err = mediaService.GetAlbum(albumID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	album.ID = uint(albumID)
	album.CreatedAt = existingAlbum.CreatedAt
	err = mediaService.UpdateAlbum(album)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, album)
}

func DeleteAlbum(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	albumID := uint64(c.GetInt64("ID"))
	album, err := mediaService.GetAlbum(albumID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	if err := mediaService.DeleteAlbum(&album); err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func GetAlbum(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	albumID := uint64(c.GetInt64("ID"))
	album, err := mediaService.GetAlbum(albumID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, album)
}

func GetAlbums(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	cursor := c.Query("cursor")
	userID := c.Query("userId")
	albums, nextCursor, err := mediaService.GetAlbums(cursor, userID)
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"nextCursor": nextCursor, "results": albums})
}

func GetTags(c *gin.Context) {
	mediaService := getMediaServiceFromContext(c)
	tags, err := mediaService.GetTags()
	if err != nil {
		handleServiceError(err, c)
		return
	}
	c.JSON(http.StatusOK, tags)
}
