package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/willdady/mediams/internal/mediams/handlers"
	"github.com/willdady/mediams/internal/mediams/models"
	"github.com/willdady/mediams/internal/mediams/postgres"
	"github.com/willdady/mediams/internal/rest"
	"github.com/willdady/mediams/internal/utils"
)

var (
	pgHost             = utils.Getenv("PG_HOST", "0.0.0.0")
	pgPort             = utils.Getenv("PG_PORT", "5432")
	pgUser             = utils.Getenv("PG_USER", "postgres")
	pgDB               = utils.Getenv("PG_DB", "postgres")
	pgPassword         = utils.Getenv("PG_PASSWORD", "mysecretpassword")
	pgSSLMode          = utils.Getenv("PG_SSL_MODE", "disable")
	dbConnectionString = fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v", pgHost, pgPort, pgUser, pgDB, pgPassword, pgSSLMode)
)

func connectToDB(retry int) (db *gorm.DB, err error) {
	if retry == 5 {
		err := errors.New("Failed to connect to database after 5 tries")
		return nil, err
	}
	db, dbErr := gorm.Open("postgres", dbConnectionString)
	if dbErr != nil {
		duration := time.Second + time.Second*time.Duration(retry)
		log.Println(dbErr)
		log.Printf("Failed to connect to database. Retrying in %v seconds.\n", duration.Seconds())
		time.Sleep(duration)
		return connectToDB(retry + 1)
	}
	return db, nil
}

var resources = rest.ResourceMap{
	"media": rest.ActionMap{
		"create": handlers.CreateMedia,
		"detail": handlers.GetMedia,
		"list":   handlers.GetMedias,
		"update": handlers.UpdateMedia,
		"delete": handlers.DeleteMedia,
	},
	"albums": rest.ActionMap{
		"create": handlers.CreateAlbum,
		"detail": handlers.GetAlbum,
		"list":   handlers.GetAlbums,
		"update": handlers.UpdateAlbum,
		"delete": handlers.DeleteAlbum,
	},
	"tags": rest.ActionMap{
		"list": handlers.GetTags,
	},
}

func main() {
	db, err := connectToDB(0)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.Album{})
	db.AutoMigrate(&models.Media{})

	mediaService := postgres.NewMediaService(db)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("mediaService", mediaService)
		c.Next()
	})

	rest.AttachEndpoints(resources, r)

	// TODO: Support PORT env variable
	r.Run() // listen and serve on 0.0.0.0:8080
}
