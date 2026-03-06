package models

type Watchlist struct {
	UserID  uint  `json:"user_id" gorm:"foreignKey:UserId"`
	MovieID uint  `json:"movie_id" gorm:"foreignKey:MovieID"` //
	Movie   Movie `json:"movie" gorm:"foreignKey:MovieID"` //
}
