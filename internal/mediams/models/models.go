package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
)

// CommonFields defines fields common to all models
type CommonFields struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
}

// Album model
type Album struct {
	CommonFields
	UserID string `json:"userId" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

// Media model
type Media struct {
	CommonFields
	AlbumID   uint           `json:"albumId"`
	UserID    string         `json:"userId" binding:"required"`
	Url       string         `json:"url" binding:"required"`
	Width     uint           `json:"width"`
	Height    uint           `json:"height"`
	MediaType string         `json:"mediaType"`
	Tags      pq.StringArray `json:"tags" gorm:"type:varchar(64)[]"`
}

// Validate validates the struct's values returning an error if invalid
func (m *Media) Validate() (err error) {
	if !govalidator.IsURL(m.Url) {
		// TODO: return a ValidationError
		return errors.New("invalid url")
	}
	return
}

// BeforeSave hook simply calls Validate
func (m *Media) BeforeSave() (err error) {
	return m.Validate()
}

// BeforeUpdate hook simply calls Validate
func (m *Media) BeforeUpdate() (err error) {
	return m.Validate()
}
