package model

type Comment struct {
	ID         int64  `gorm:"primaryKey;type:uint" json:"id"`
	VideoID    int64  `gorm:"type:uint" json:"video_id"`
	UserID     int64  `gorm:"type:mediumint unsigned" `
	UserMsg    User   `gorm:"foreignKey:UserID" json:"user"`
	CommentMsg string `gorm:"type:varchar(255)" json:"content"`
	CreateDate string `gorm:"type:varchar(20)" json:"create_date"` //评论发布日期，格式 yy-mm-dd
	DeleteDate string `gorm:"type:varchar(20)" json:"delete_date"` //评论软删除
}

func (table *Comment) TableName() string {
	return "comment" // 指定表名，大小写敏感
}
