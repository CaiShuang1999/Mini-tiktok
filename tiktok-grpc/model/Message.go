package model

type Message struct {
	ID         int64  `gorm:"primaryKey;type:uint" json:"id"`
	ToUserID   int64  `gorm:"type:mediumint unsigned" json:"to_user_id"`
	FromUserID int64  `gorm:"type:mediumint unsigned" json:"from_user_id"`
	Content    string `gorm:"type:varchar(255)" json:"content"`
	CreateTime int64  `gorm:"type:uint" json:"create_time"`
}

func (table *Message) TableName() string {
	return "message" // 指定表名，大小写敏感
}
