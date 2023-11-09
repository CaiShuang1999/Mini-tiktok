package model

type User struct {
	ID              int64  `gorm:"primaryKey;type:mediumint unsigned" json:"id"  redis:"-"`                            // 用户id(gorm主键自增)
	Name            string `gorm:"not null;unique;size:32" json:"name" redis:"name"`                                   // 用户名称
	Password        string `gorm:"not null" json:"-" redis:"-"`                                                        //密码(返回json时，不包括密码)
	Avatar          string `gorm:"default:'/default/avatar.jpg'" json:"avatar" redis:"avatar"`                         // 用户头像
	BackgroundImage string `gorm:"default:'/default/background.jpg'" json:"background_image" redis:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `gorm:"type:mediumint unsigned" json:"favorite_count" redis:"favorite_count"`               // 喜欢数
	WorkCount       int64  `gorm:"type:smallint(5) unsigned" json:"work_count" redis:"work_count"`                     // 作品数
	FollowCount     int64  `gorm:"type:smallint(5) unsigned" json:"follow_count" redis:"follow_count"`                 // 关注总数
	FollowerCount   int64  `gorm:"type:mediumint unsigned" json:"follower_count" redis:"follower_count"`               // 粉丝总数
	IsFollow        bool   `gorm:"-" json:"is_follow" redis:"-"`                                                       // true-已关注，false-未关注
	Signature       string `gorm:"default:'默认个人签名!'" json:"signature,omitempty"  redis:"signature"`                    // 个人简介
	TotalFavorited  int64  `gorm:"type:mediumint unsigned" json:"total_favorited" redis:"-"`                           // 获赞数量

}

func (table *User) TableName() string {
	return "user" // 指定表名，大小写敏感
}
