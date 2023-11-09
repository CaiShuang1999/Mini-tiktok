package model

type Relation struct {
	UserID   int64 `gorm:"primaryKey;type:mediumint unsigned;autoIncrement:false"`
	ToUserID int64 `gorm:"primaryKey;type:mediumint unsigned;autoIncrement:false"`
	IsFollow bool  `json:"is_follow"`
}

func (table *Relation) TableName() string {
	return "relation" // 指定表名，大小写敏感
}
