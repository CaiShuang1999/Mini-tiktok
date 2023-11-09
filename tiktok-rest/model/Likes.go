package model

type Favorite struct {
	VideoId    int64 `gorm:"primaryKey;type:uint;autoIncrement:false" json:"video_id"`
	UserID     int64 `gorm:"primaryKey;type:mediumint unsigned;autoIncrement:false"`
	IsFavorite bool  `json:"is_favorite"`
}

func (table *Favorite) TableName() string {
	return "favorite" // 指定表名，大小写敏感
}
