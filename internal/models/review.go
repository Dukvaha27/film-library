package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	MovieID uint   `gorm:"not null;index" json:"movie_id"`
	Author  string `gorm:"size:100;not null" json:"author"`
	Rating  int    `gorm:"not null" json:"rating"`
	Comment string `gorm:"type:text" json:"comment"`

	Movie Movie `gorm:"foreignKey:MovieID" json:"-"`
}
