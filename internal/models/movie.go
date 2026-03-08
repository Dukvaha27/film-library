package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model

	AvgRating float64  `gorm:"type:decimal(3,2);default:0" json:"avg_rating"`
	Reviews   []Review `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE;" json:"reviews,omitempty"`
	Genres    []Genre  `gorm:"many2many:movie_genres;" json:"genres,omitempty"`
}
