package model

type Video struct {
	ID            int64  `gorm:"primaryKey type:uint" json:"id"  redis:"-"` // 视频唯一标识
	UserID        int64  `gorm:"type:mediumint unsigned"  redis:"user_id"`
	Author        User   `gorm:"foreignKey:UserID" json:"author"  redis:"-"`            // 视频作者信息
	CommentCount  int64  `gorm:"type:uint" json:"comment_count"  redis:"comment_count"` // 视频的评论总数
	CoverURL      string `json:"cover_url"  redis:"cover_url"`                          // 视频封面地址
	FavoriteCount int64  `gorm:"type:uint" json:"favorite_count" redis:"-"`             // 视频的点赞总数
	IsFavorite    bool   `gorm:"-" json:"is_favorite" redis:"-"`                        // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url" redis:"play_url"`                             // 视频播放地址
	Title         string `gorm:"default:'默认标题'" json:"title"  redis:"title"`            // 视频标题
	CreateTime    int64  `gorm:"type:uint" json:"-"  redis:"create_time"`
}

func (table *Video) TableName() string {
	return "video" // 指定表名，大小写敏感
}
