package models

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
	Name string `gorm:"size:100;not null;unique" json:"name"`

	Movies []Movie `gorm:"many2many:movie_genres;" json:"movies,omitempty"`
}
