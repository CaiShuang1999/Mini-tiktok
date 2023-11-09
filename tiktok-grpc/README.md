
# Mini-tiktok

字节跳动后端青训-迷你抖音app后端项目(gin+grpc微服务实现)

api文档：https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707523

### 技术栈
gin+gRPC+ etcd+gorm+redis+viper+JWT+Nginx+ffmpeg+zap


### 目录
 [微服务](#微服务)
> -  [基础概念](#基本概念)    
> - [微服务数据传输](#微服务数据传输) 
> - [微服务划分](#微服务划分) 
> - [文件目录](#文件) 
> - [服务发现与注册](#服务发现与注册) 
- [静态文件配置](#静态文件服务)

- [model](#model)

[微服务逻辑实现](#service) 
> - [user](#user)     
> - [video](#video) 
> - [feed](#feed) 
> - [favorite](#favorite) 
>  - [comment](#comment) 
> - [relation](#relation) 
> - [message](#message)

- [缓存设计](#缓存设计)

  
## 微服务

###  基本概念
### Restful api 
- 一切数据视作资源，并通过 URI（统一资源标识符）来标识和访问这些资源
- HTTP 方法（GET、POST、PUT、DELETE 等）来表示对资源的不同操作（增删改查）
- HTTP响应状态码，描述操作的结果
### RPC:远程过程调用 
 概念：
 - 构建分布式应用程序的通信模式。它允许一个计算机程序通过网络调用另一个计算机程序中的函数或方法，就像调用本地函数一样。
- 在传统的过程调用中，函数调用是在同一个进程内进行的，而在 RPC 中，函数调用是在不同的计算机或进程之间进行的，通常是通过网络。

微服务应用：
- 微服务架构是一种将复杂应用程序拆分成一系列小型、独立的服务的架构风格。每个微服务都可以独立开发、部署和扩展，通过轻量级的通信方式进行交互。
- RPC 提供了一种方便的方式来实现微服务之间的通信。每个微服务可以通过 RPC 调用其他微服务的函数或方法，以完成特定的业务逻辑。
  

### 微服务数据传输：
注意： 当 gRPC 服务端传输 Protocol Buffers 消息时，如果字段的值为默认值，该字段在传输过程中是不会被包含的。这是为了节省带宽和提高效率。（但在客户端接收得到默认值res.StatusCode=0）

---------手动在生成的pb文件响应中删除json的omitempty的tag-----------------------
> 请求：gin(http:json) --->  rpc 客户端（tcp:json转为message）--->  rpc 客户端（tcp:message）
  响应：rpc 客户端（tcp:message）--->  rpc 客户端（tcp:message转为json）

### 微服务划分：
user：用户注册，登录，用户信息
feed：视频流
video：上传视频，视频信息
favorite：点赞(取消赞)，点赞列表
comment：评论(删除评论)，评论列表
relation:关注(取关)，关注列表，好友列表
message：发送消息，消息列表


### 文件
```
/
|-- apps        //微服务
|   |-- user
|   |	|-- api  // 定义路由
|   |   |-- pb   //proto文件生成
|   |   |-- rpc   
|   |   | 	|-- client  // 客户端（作为gin与rpc服务端中间）
|   |   | 	|-- server  //服务端
|   |   | 	|   |-- service.go  // 具体服务代码实现
|   |   | 	|   |-- rpcServerInit.go // rpc服务端启动代码
|   |   |-- user.proto //定义proto文件  
|   |-- video
|   |-- feed
|   |-- favorite
|   |-- relation
|   |-- comment
|   |-- message
|   |-- comment
|   |--init_services.go  //协程启动所有微服务

|-- cmd
|   |-- init_database.go   //sql启动
|   |-- init_nginx.go   //nginx
|   |-- init_redis.go   //redis启动
|   |-- init_router.go   //所有微服务router启动
|-- common
|   |-- jwtx   //jwt
|   |-- redisx  //redis缓存代码
|   |-- zapx  //日志
|   |-- utils  
|   |   |-- convert_msg.go   //消息与模型的转换(gorm数据库记录绑定到模型，但rpc传输的是pb.message就要进行类型转换)
|   |   |-- hash.go   //用户密码加密
|-- configs
|-- model
|-- web   
|   |-- static   //静态文件(上传的视频)
|-- main.go  //主入口
```


 ### 服务发现与注册
为什么需要服务注册与发现
>在使用微服务架构时，我们会将一个大的单应用拆分成多个独立自治的小服务，如果在没有服务发现的情况下，我们要想在服务之间进行通信，我们只能使用硬编码的方式，将需要通信的服务配置信息写在程序中，这样可能会导致一系列的问题，如服务提供方的网络发生了变化，服务的调用方如不及时修改将会影响使用，无法动态收缩和扩容等。

  


服务注册与发现中心的职责
>1.  管理当前注册到服务注册与发现中心的微服务实例元数据信息，包括服务实例的服务名、IP地址、端口号、服务状态和服务描述信息等。
>2.  与注册到服务注册与发现中心的微服务实例维持心跳，定期检查注册表中的服务实例是否在线，并剔除无效服务实例信息。
>3.  提供服务发现能力，为调用方提供服务方的服务实例元数据。





  



  
  
  
  

## 静态文件服务
### (gin或Nginx)

  

assets：静态文件目录(default下面存放默认的头像和背景图,videos保存上传的视频，video_covers视频封面）

  

1.使用gin的Static

r := gin.Default()

r.Static("/static", "./assets")

服务器IP地址:8080/static->映射实际assets目录

  
  

2.使用Nginx进行配置(listen：服务器(IP:80) ；root:实际静态文件目录的路径)

  
  

server {

listen localhost:80;

location / {

root E:/GoWeb/assets;

autoindex on;

}

客户端请求默认头像：

(服务器IP地址:80/default/avatar.jpg)->root/default/avatar.jpg

  

## model

根据api文档response所有必需的Json字段都需满足(字符串不能为空)

1. user：

> Avatar string `gorm:"default:'/default/avatar.jpg'" json:"avatar"` // 用户头像

  

>BackgroundImage string `gorm:"default:'/default/background.jpg'" json:"background_image"` // 用户背景

  

response返回的user包括了关注关系,但用 gorm:"-" 不把此字段存入user表中

>IsFollow bool `gorm:"-" json:"is_follow"`

  

response不返回密码,用 json:"-"

>Password string `gorm:"not null;size:32" json:"-"`

  

2. video：

video表中只有user_id字段通过UserID外键关联user表中作者信息

> UserID int64

> Author User `gorm:"foreignKey:UserID" json:"author"` // 视频作者信息

  

response返回的video包括了点赞关系, gorm:"-" 不把此字段存入video表中

>IsFavorite bool `gorm:"-" json:"is_favorite"` // true-已点赞，false-未点赞

  

3. favorite：关于点赞表设计: （video_id,user_id）

- 方案1.点赞数据表就加一条记录,取消赞就删除这条记录

>  - 查询某个用户是否点赞时只需要查找该用户的点赞记录，没有点赞就表示没有记录。这种查询性能相对较好，因为只需要进行一次查询操作。

>  - 每次点赞都是插入一条新记录

>  - 取消赞时需要删除相应的记录，这可能会导致表的性能下降，特别是在点赞数据量较大时

- 方案2.设置一个字段status,点赞就加一条记录,status变成1, 取消赞 就把status变成0

>  - 查询点赞关系，需要查询点赞表并检查状态字段的值。

>  - 点赞和取消赞时只需要更新一个字段的值，性能方面，更有优势。

  

方案一简单点，方案二在面对高数据量性能更好

  

4. relation：关于关注表设计: （user_id,to_user_id） ：同点赞表选择方案2

  

5. comment:

- 通过userID外键查看发送者信息

>UserID int64

UserMsg User `gorm:"foreignKey:UserID" json:"user"`

DeleteDate string `json:"delete_date"` //评论软删除(空值为未删除)

6. message:

>CreateTime int64 `json:"create_time"` 时间戳

  
  

## Service
### user        
#### apps/user/rpc/service.go:

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
- 判断token是否有效
- 判断用户信息是否在缓存(有：直接取；无：从数据库中取，再缓存)
- 返回响应(user信息)

### video:

#### apps/video/rpc/service.go:

1.上传视频：POST/douyin/publish/action/
>Body 参数 (application/form-data)：
>- data(*multipart.FileHeader)视频数据
>- Token
>- Title
- 判断token是否过期
- 上传视频到服务器静态资源目录
- ffmpeg截取封面上传到服务器静态资源目录（文件名命名格式：id+上传时间戳(只能单次上传一个文件)）
- 上传成功，video表添加记录，user表中work_count加1（这里修改了user表，为了保证数据库与缓存的一致性，把缓存的user信息删除）
- 将视频信息缓存
- 返回基础响应

  

2.用户发布列表：GET/douyin/publish/list/
>query(user_id,token)

- 不判断token(后面如果通过feed流查看作者，如果使用了token验证会看不见别人的发布列表)
- 通过user_id==video.user_id找到用户上传的所有视频（db.Preload("Author").Where("user_id = ?", userID).Find(&videos) 通过userID外键连接的"Author"需要.Preload方法预加载，否则"Author"为空值）（找到的videos还要进行类型转化为message类型）
- 所有视频videos还有判断作者的关注关系，和点赞关系
- 返回响应（videos）

  
### feed:

#### apps/feed/rpc/service.go:
视频流：GET/douyin/feed/ 返回最新的最多30个投稿
>query(latest_time,token)
- token判断
>- token为空未登录状态，只能刷视频无法（关注，点赞，评论）
>- token不为空，判断是否过期
- 根据id倒序找到最新的videos
- 判断user与当前video的点赞关系和作者的关注关系(查表，如果有点赞记录且is_favorite为true，则设置返回响应的video的is_favorite为true；如果关注了作者，则"author"中的"is_follow"设为true)
- 返回响应（videos）

  

### favorite:
#### apps/favorite/rpc/service.go:
1.赞操作：POST/douyin/favorite/action/
>query(token,video_id,action_type)
- token判断是否过期
- 通过token的user_id和video_id进行查找点赞表是否存在记录，如果不存在则新建(每次操作都得先查询是否存在记录)
>- 根据action_type，如果是点赞操作，则将记录的"is_favorite"设为true
>- 如果是取消赞操作，则将记录的"is_favorite"设为false
>- 点赞(取消赞)操作后将user表中自己的喜欢视频数"favorite_count"+1(-1),user表中视频作者的总获赞数"total_favorited"+1(-1),video表中该视频的点赞数"favorite_count"+1(-1)
- 删除user:id,video:id缓存,更新favorite:user:id，favorite:video:id缓存
- 返回基本响应

  
2.喜欢列表：GET/douyin/favorite/list/
>query(token,user_id)
- token判断是否过期
- 通过favorite.user_id = user_id AND favorite.is_favorite=true找到用户点赞过的所有视频再连接user表返回视频作者信息(预加载)
- videos每个视频的"is_favorite"设为true
- 返回响应(videos)



### comment:
#### apps/comment/rpc/service.go:
1.评论操作：POST/douyin/comment/action/

>query(token,video_id,action_type,commnet_text,comment_id)
- token判断是否过期
- 通过action_type进行操作匹配（1-发布评论，2-删除评论）
> -评论：创建comment记录(token.user_id,video_id,comment_text,create_time)
>  - 删除：通过id找到评论，判断评论的用户与当前用户是否是同一个id（无法删除他人评论）；采用软删除(通过评论id，填充delete_time)
-返回基础响应

2.评论列表：GET/douyin/comment/list/
>query(token,video_id)
- token判断是否过期
- 通过video_id找到delete_time为空(未删除)的所有评论再以评论id倒序返回comments(预加载评论user)
- 返回响应(comments)

  

### relation.go:
#### apps/relation/rpc/service.go:
1.关注操作：POST/douyin/relation/action/
>query(token,to_user_id,action_type)
- token判断是否过期
- 判断to_user_id与当前用户是否为同一个(不能关注自己)
- 通过token的user_id和to_user_id进行查找关系表是否存在记录，如果不存在则新建(每次操作都得先查询是否存在记录)
>  * 关注：修改"is_follow"为true
>  * 取消关注：修改"is_follow"为false
>  * 关注(取关)当前用户user表"follow_count"数+1(-1)，被关注用户to_user_id的"follower_count"+1(-1)
- 删除用户缓存信息
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

  

### message:只能和好友（互关）操作
#### apps/message/rpc/service.go:
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

- 先更新 db
- 然后直接删除 cache 。

**读** :

- 从 cache 中读取数据，读取到就直接返回
- cache 中读取不到的话，就从 db 中读取数据返回
- 再把数据放到 cache 中。

  
 
#### Write Behind Pattern（异步缓存写入）：

**只更新缓存，不直接更新 db，而是改为异步批量的方式来更新 db**
这种方式对数据一致性带来了更大的挑战，比如 cache 数据可能还没异步更新 db 的话，cache 服务可能就就挂掉了
(写性能非常高，非常适合一些数据经常变化又对数据一致性要求没那么高的场景，比如浏览量、点赞量。)

  
  
  
  

### 本项目：

#### Redis :

#### 旁路缓存模式

用户信息(登录操作)
视频信息(上传，feed)
> feed 需要进行排序，用有序集合类型（但随着视频集合元素增多，需要考虑大key问题）


#### 异步缓存写入：
用户：收到的赞数量
视频：点赞数量

注意：更新操作的过期时间选择是否更新过期时间(`Get` 命令获取键的当前过期时间，然后再使用 `Set` 命令更新键的值并指定之前获取的过期时间作为参数)