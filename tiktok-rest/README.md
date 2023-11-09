# Mini-tiktok
字节跳动后端青训-抖音app

一个简单的demo
api文档：https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707523

### Restful api

-  一切数据视作资源，并通过 URI（统一资源标识符）来标识和访问这些资源
- HTTP 方法（GET、POST、PUT、DELETE 等）来表示对资源的不同操作（增删改查）
- HTTP响应状态码，描述操作的结果

### 技术栈

gin+	gorm+redis+viper+JWT+Nginx+ffmpeg

### 文件


assets:	静态文件
config:数据库配置
model：数据结构
route：路由
service：服务

### 目录

- [静态文件配置](##静态文件服务(gin或Nginx))
- [model](#model)
- [service](#service) : [user](#user.go) [video](#video.go)   [favorite](#favorite.go)  [comment](#comment.go)   [relation](#relation.go) [message](#message.go) 
- [缓存设计](#缓存设计)




## 静态文件服务(gin或Nginx)

assets：静态文件目录(default下面存放默认的头像和背景图,videos保存上传的视频，video_covers视频封面）

1.使用gin的Static
r  := gin.Default()
r.Static("/static", "./assets")
服务器IP地址:8080/static->映射实际assets目录
  


2.使用Nginx进行配置(listen：服务器(IP:80)	；root:实际静态文件目录的路径)


	server {
        listen      localhost:80;
        location / {
           root   E:/GoWeb/assets;
			autoindex on;
			}
   
 客户端请求默认头像：
 (服务器IP地址:80/default/avatar.jpg)->root/default/avatar.jpg

## model
根据api文档response所有必需的Json字段都需满足(字符串不能为空)
1. user：
> Avatar string  `gorm:"default:'/default/avatar.jpg'" json:"avatar"`  // 用户头像

>BackgroundImage string  `gorm:"default:'/default/background.jpg'" json:"background_image"` // 用户背景

response返回的user包括了关注关系,但用 gorm:"-" 不把此字段存入user表中
>IsFollow bool  `gorm:"-" json:"is_follow"`

response不返回密码,用 json:"-" 
>Password string  `gorm:"not null;size:32" json:"-"`

2. video：
video表中只有user_id字段通过UserID外键关联user表中作者信息
> UserID int64
> Author User `gorm:"foreignKey:UserID" json:"author"`  // 视频作者信息

response返回的video包括了点赞关系, gorm:"-" 不把此字段存入video表中
>IsFavorite bool  `gorm:"-" json:"is_favorite"`  // true-已点赞，false-未点赞

3. favorite：关于点赞表设计: （video_id,user_id） 
- 方案1.点赞数据表就加一条记录,取消赞就删除这条记录
> - 查询某个用户是否点赞时只需要查找该用户的点赞记录，没有点赞就表示没有记录。这种查询性能相对较好，因为只需要进行一次查询操作。
> - 每次点赞都是插入一条新记录
> - 取消赞时需要删除相应的记录，这可能会导致表的性能下降，特别是在点赞数据量较大时
- 方案2.设置一个字段status,点赞就加一条记录,status变成1, 取消赞 就把status变成0
> - 查询点赞关系，需要查询点赞表并检查状态字段的值。
> - 点赞和取消赞时只需要更新一个字段的值，性能方面，更有优势。

方案一简单点，方案二在面对高数据量性能更好

4. relation：关于关注表设计: （user_id,to_user_id） ：同点赞表选择方案2

5. comment:
- 通过userID外键查看发送者信息
>UserID int64
UserMsg User `gorm:"foreignKey:UserID" json:"user"`
DeleteDate string  `json:"delete_date"`  //评论软删除(空值为未删除)
6. message:
>CreateTime int64  `json:"create_time"` 时间戳


## Service
### user.go:
1.注册：POST/douyin/user/register/                 
>query(user_name,password)
- 用户名和密码长度校验（小于32个字符）
- 用户名唯一校验
- bcrypt 密码哈希函数，安全的密码存储方法
- 注册成功，生成token；返回响应(user_id,token)

2.登录：POST/douyin/user/login/              
>query(user_name,password)
- 用户名和密码长度校验（小于32个字符）
- 用户名和密码匹配
- 登录成功，生成token；返回响应(user_id,token)

3.用户界面：GET/douyin/user/       
>query(user_id,token)
- 判断token是否有效(token是否过期，query的user_id是否等于token里的user_id)
- 判断用户信息是否在缓存(有：直接取；无：从数据库中取，再缓存)
- 返回响应(user信息)


### video.go:

1.上传视频：POST/douyin/publish/action/     
>Body 参数 (application/form-data)：
>- data(*multipart.FileHeader)视频数据 
>- Token
>- Title

- 判断token是否过期
- 返回响应(user信息)
- c.SaveUploadedFile(file, savePath)上传视频到服务器静态资源目录
- ffmpeg截取封面上传到服务器静态资源目录（文件名命名格式：id+上传时间戳(只能单次上传一个文件)）
- 上传成功，video表添加记录，user表中work_count加1
- 将视频缓存
- 返回基础响应

2.用户发布列表：GET/douyin/publish/list/   
>query(user_id,token)

- 不判断token(后面如果通过feed流查看作者，如果使用了token验证会看不见别人的发布列表)
- 通过user_id==video.user_id找到用户上传的所有视频（db.Preload("Author").Where("user_id = ?", userID).Find(&videos) 通过userID外键连接的"Author"需要.Preload方法预加载，否则"Author"为空值）
- 返回响应（videos）

3.视频流：GET/douyin/feed/ 返回最新的最多30个投稿
>query(latest_time,token)

- token判断
>- token为空未登录状态，只能刷视频无法（关注，点赞，评论）
>- token不为空，判断是否过期
- 根据id倒序找到最新的videos
- 判断user与当前video的点赞关系和作者的关注关系(查表，如果有点赞记录且is_favorite为true，则设置返回响应的video的is_favorite为true；如果关注了作者，则"author"中的"is_follow"设为true)
- 返回响应（videos）

### favorite.go:
1.赞操作：POST/douyin/favorite/action/
>query(token,video_id,action_type)

- token判断是否过期
- 通过token的user_id和video_id进行查找点赞表是否存在记录，如果不存在则新建(每次操作都得先查询是否存在记录)
>- 根据action_type，如果是点赞操作，则将记录的"is_favorite"设为true
>- 如果是取消赞操作，则将记录的"is_favorite"设为false
>- 点赞(取消赞)操作后将user表中自己的喜欢视频数"favorite_count"+1(-1),user表中视频作者的总获赞数"total_favorited"+1(-1),video表中该视频的点赞数"favorite_count"+1(-1)

- 返回基本响应

2.喜欢列表：GET/douyin/favorite/list/
>query(token,user_id)

- token判断是否过期
- 通过favorite.user_id = user_id AND favorite.is_favorite=true找到用户点赞过的所有视频再连接user表返回视频作者信息(预加载)
- videos每个视频的"is_favorite"设为true
- 返回响应(videos)
- 
### comment.go:
1.评论操作：POST/douyin/comment/action/
>query(token,video_id,action_type,commnet_text,comment_id)

- token判断是否过期
- 通过action_type进行操作匹配（1-发布评论，2-删除评论）
> -评论：创建comment记录(token.user_id,video_id,comment_text,create_time)
> - 删除：通过id找到评论，判断评论的用户与当前用户是否是同一个id（无法删除他人评论）；采用软删除(通过评论id，填充delete_time)
-返回基础响应

2.评论列表：GET/douyin/comment/list/
>query(token,video_id)

- token判断是否过期
- 通过video_id找到delete_time为空(未删除)的所有评论再以评论id倒序返回comments(预加载评论user)
- 返回响应(comments)

### relation.go:

1.关注操作：POST/douyin/relation/action/
>query(token,to_user_id,action_type)

- token判断是否过期
- 判断to_user_id与当前用户是否为同一个(不能关注自己)
- 通过token的user_id和to_user_id进行查找关系表是否存在记录，如果不存在则新建(每次操作都得先查询是否存在记录)
> * 关注：修改"is_follow"为true
>  * 取消关注：修改"is_follow"为false
>  * 关注(取关)当前用户user表"follow_count"数+1(-1)，被关注用户to_user_id的"follower_count"+1(-1)
- 返回基础响应

2.关注列表：GET/douyin/relation/follow/list/
>query(token,user_id)
- token判断是否过期
- 通过当前用户ID(relation.user_id)和"is_follow=true"查找关系表，找到当前用户关注的所有用户ID并通过连接user表找出用户信息follows
- 设定follows"is_follow"为true
- 返回响应follows

3.粉丝列表：GET/douyin/relation/follower/list/
>query(token,user_id)
- token判断是否过期
- 通过当前用户ID(relation.to_user_id)和"is_follow=true"查找关系表，找到关注了当前用户的所有用户ID并通过连接user表找出用户信息fans
- 再查找关注表找到当前用户是否关注了他的粉丝(如果互相关注设定fans的"is_follow"为true)
- 返回响应fans

4.好友列表：GET/douyin/relation/friend/list/
>query(token,user_id)
- token判断是否过期
- 找到当前用户（user_id）相互关注的用户(user)
- 设定follows"is_follow"为true
- 返回响应friends

### message.go:只能和好友（互关）操作
1.发送消息：POST/douyin/message/action/
>query(token,to_user_id,action_type,content)
- token判断是否过期
- action_type(1-发送消息)创建message记录
- 返回基础响应

2.聊天记录：POST/douyin/message/action/
>query(token,to_user_id,pre_msg_time)
- token判断是否过期
- 采用轮询(每秒)pre_msg_time代表最新一条消息的时间戳
- 查找message表(当前用户发送给to_user_id的消息和to_user_id发送给当前用户的消息)且"create_time">pre_msg_time()
- 返回响应messages

## 缓存设计

### Redis

缓存数据一致性：https://xiaolincoding.com/redis/architecture/mysql_redis_consistency.html#%E6%95%B0%E6%8D%AE%E5%BA%93%E5%92%8C%E7%BC%93%E5%AD%98%E5%A6%82%E4%BD%95%E4%BF%9D%E8%AF%81%E4%B8%80%E8%87%B4%E6%80%A7

https://javaguide.cn/database/redis/3-commonly-used-cache-read-and-write-strategies.html#cache-aside-pattern-%E6%97%81%E8%B7%AF%E7%BC%93%E5%AD%98%E6%A8%A1%E5%BC%8F

https://www.bilibili.com/video/BV1Dy4y1976U/?spm_id_from=333.337.search-card.all.click&vd_source=5fd2598c2a7f4ad650ee91f5d50dc8f1


#### Cache Aside Pattern（旁路缓存模式）：比较适合读请求比较多的场景
**写**：

-   先更新 db
-   然后直接删除 cache 。

**读**  :

-   从 cache 中读取数据，读取到就直接返回
-   cache 中读取不到的话，就从 db 中读取数据返回
-   再把数据放到 cache 中。



#### Write Behind Pattern（异步缓存写入）：

**只更新缓存，不直接更新 db，而是改为异步批量的方式来更新 db**
这种方式对数据一致性带来了更大的挑战，比如 cache 数据可能还没异步更新 db 的话，cache 服务可能就就挂掉了

(写性能非常高，非常适合一些数据经常变化又对数据一致性要求没那么高的场景，比如浏览量、点赞量。)




### 本项目：
#### Redis :
#### 旁路缓存模式
用户信息(登录操作)
视频信息(上传，feed)
，评论，消息

#### 异步缓存写入：
用户：收到的赞数量
视频：点赞数量


注意：更新操作的过期时间选择是否更新过期时间(`Get` 命令获取键的当前过期时间，然后再使用 `Set` 命令更新键的值并指定之前获取的过期时间作为参数)