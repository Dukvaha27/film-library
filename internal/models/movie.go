package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title       string  `gorm:"size:200;not null" json:"title" validate:"required"`
	Description string  `gorm:"type:text" json:"description"`
	Year        int     `gorm:"index" json:"year"`
	Duration    int     `json:"duration"`
	GenreID     uint    `gorm:"index" json:"genre_id" validate:"required"`
	Genre       Genre   `gorm:"foreignKey:GenreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	AvgRating   float32 `gorm:"default:0" json:"avg_rating"`
	ViewsCount  int     `gorm:"default:0" json:"views_count"`
}

type UpdateMovie struct {
	Title       string  `gorm:"size:200;not null" json:"title,omitempty"`
	Description string  `gorm:"type:text" json:"description,omitempty"`
	Year        int     `gorm:"index" json:"year,omitempty"`
	Duration    int     `json:"duration,omitempty"`
	AvgRating   float32 `json:"avg_rating,omitempty"`
	ViewsCount  int     `json:"views_count,omitempty"`
}
