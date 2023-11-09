
# Mini-tiktok

字节跳动后端青训-迷你抖音app后端项目(gin+grpc微服务实现)

api文档：https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707523

### 技术栈
go-zero+ etcd+gorm+redis+viper+JWT+Nginx+ffmpeg+zap


### 目录

 - [文件目录](#文件目录) 
- [静态文件配置](#静态文件服务)

 [Go-zero框架](#go-zero框架)
> -  [API](#api)    
> - [RPC](#rpc) 

[微服务](#微服务)
> - [微服务划分](#微服务划分) 
> - [服务发现与注册](#服务发现与注册) 
> - [负载均衡](#负载均衡) 


[微服务逻辑实现](#微服务逻辑实现) 
> - [user](#user)     
> - [video](#video) 
> - [feed](#feed) 
> - [favorite](#favorite) 
>  - [comment](#comment) 
> - [relation](#relation) 
> - [message](#message)

- [缓存设计](#缓存设计)

  
## 文件目录
```
| -- tiktok-micro
| |-- app
| | |-- api  	-API
| | |-- services -微服务
| |-- common
| | |-- cryptx  	-密码哈希加密
| | |-- jwtx 
| |-- model
| |-- web
| | |-- static -静态服务地址(nginx)
| | | |-- image -图片(头像，背景，视频封面)
| | | |-- video -用户上传的视频


```
##  静态文件服务
使用Nginx进行配置(listen：服务器(IP:80) ；root:实际静态文件目录的路径)
server {
listen localhost:80;
location / {
root web/static;
autoindex on;
}

 
客户端请求默认头像：
(服务器IP地址:80/image/avatar.jpg)->root/image/avatar.jpg



##  Go-zero框架
###  API
API 网关作为一个中间层，负责处理请求的路由和转发，将请求分发给各个微服务进行处理，并将响应返回给客户端
#### 定义app.api
1.api文件定义：
- 请求（query格式）
- 响应(JSON格式)
- 调用方法
- 服务分组 (https://go-zero.dev/docs/tutorials/api/route/group)

2.生成api： goctl api go -api app.api -dir ./api -m (-m：服务分组)

```
/

| -- api
| |-- etc 
| | |-- app.yml  - 配置定义(接口监听,微服务rpc地址)
| |-- internal 
| | |-- config  
| | |-- handler
| | |  |-- router.go  //路由与中间件配置
| | |  |-- ...handler.go  //handler(解析请求传给logic，返回logic传回的响应)
| | |-- logic  //微服务RPC请求响应的类型转换，响应返回handler
| | |-- svc  
| | |-- types -根据api文件生成的请求和响应类型定义
```

### RPC
在services文件下创建各个微服务：如user

1.proto文件定义：
- 包名
> option go_package="./user_pb";
package  user_pb;
- 请求（query格式）
- 响应(JSON格式)
- 调用方法


2.生成rpc文件:goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.   
```
/

| -- user
| |-- etc 
| | |-- user.yml  - 配置定义(微服务地址,数据库和其他)
| |-- internal 
| | |-- config
| | |-- logic  //具体业务代码
| | |-- server  
| | |-- svc  
| |-- user_pb   //proto生成的响应和方法包
| |-- userservice
| | |-- user.proto //proto文件
```
3.流程
(1)配置：
 - (etc/user.yml)  配置文件
- (internal/config/config.go)
 > yml文件的配置绑定为结构体

- (internal/svc/servicecontext.go)
 >   注册服务上下文 的依赖(共享的上下文，适合设置一些全局可用的资源,如配置数据库)
 >   type ServiceContext

4.微服务之间的RPC调用（如video调用user微服务获取作者信息）

-  (etc/video.yml)  配置文件添加用户rpc服务地址
>UserRpc:
    Endpoints:
 > 用户rpc端口


- internal/svc/servicecontext.go 文件中添加 user 客户端连接
>return  &ServiceContext{
Config: c,
UserRpc: userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)), // user微服务连接
}

之后就能在logic代码中通过l.svcCtx.UserRPC调用获取用户信息的方法




## 微服务


### 微服务划分
user：用户注册，登录，用户信息
feed：视频流
video：上传视频，作品列表，视频信息
favorite：点赞(取消赞)，点赞列表，点赞状态
comment：评论(删除评论)，评论列表
relation:关注(取关)，关注列表，好友列表，关注状态
message：发送消息，消息列表



### 服务发现与注册

Go-zero使用Etcd进行服务注册与发现。

RPC微服务启动时进行服务注册
- user.yml
> Etcd:
Hosts:
-127.0.0.1:2379
Key: user.rpc

API接口配置
- app.yml
>UserRpc:
		Etcd:
		Hosts:
		-127.0.0.1:2379
		Key: user.rpc

使用etcd查看
(key / etcd 的租户 id)
>user.rpc/7587874527129894460
127.0.0.1:8081
user.rpc/7587874527129894463
127.0.0.1:8083

在服务发现时，也会通过 `etcdctl get --prefix` 指令去获取所有可用的 ip 节点

### 负载均衡 
go-zero默认使用P2C（Power of Two Choices）

P2C算法的核心思想是：从所有可用节点中随机选择两个节点，然后根据这两个节点的负载情况选择一个负载较小的节点。这样做的好处在于，如果只随机选择一个节点，可能会选择到负载较高的节点，从而导致负载不均衡；而选择两个节点，则可以进行比较，从而避免最劣选择。



## 微服务逻辑实现

### user
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

  

3.用户信息：GET/douyin/user/ （被其他微服务所调用获取用户信息方法）

>query(user_id,token)

- 判断token是否有效

- 判断用户信息是否在缓存(有：直接取；无：从数据库中取，再缓存)

- 返回响应(user信息)

### video 
  
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

3.视频信息（给其他微服务调用的获取视频信息方法）(需要调用userRPC)

  >query(video_id,token)
  
  - 通过videocache（先找缓存，再找数据库）查找视频信息(和作者author_id)
  - 调用userRPC方法通过author_id获取视频作者信息
  - 返回完整的视频信息
  
  
### favorite:
  
1.赞操作：POST/douyin/favorite/action/

>query(token,video_id,action_type)

- token判断是否过期

- 通过token的user_id和video_id进行查找点赞表是否存在记录，如果不存在则新建(每次操作都得先查询是否存在记录)

>- 根据action_type，如果是点赞操作，则将记录的"is_favorite"设为true

>- 如果是取消赞操作，则将记录的"is_favorite"设为false

>- 点赞(取消赞)操作后将user表中自己的喜欢视频数"favorite_count"+1(-1),user表中视频作者的总获赞数"total_favorited"+1(-1),video表中该视频的点赞数"favorite_count"+1(-1)

- 删除user:id,video:id缓存,更新favorite:user:id，favorite:video:id缓存

- 返回基本响应

  

2.喜欢列表：GET/douyin/favorite/list/(需要调用videoRPC)

>query(token,user_id)

- token判断是否过期

- 查找favorite数据库记录找到用户点赞的所有视频ID
- 调用videoRPC方法通过视频ID获取视频信息
- 返回响应(videos)


### comment:

1.评论操作：POST/douyin/comment/action/

  

>query(token,video_id,action_type,commnet_text,comment_id)

- token判断是否过期

- 通过action_type进行操作匹配（1-发布评论，2-删除评论）

> -评论：创建comment记录(token.user_id,video_id,comment_text,create_time)

>  - 删除：通过id找到评论，判断评论的用户与当前用户是否是同一个id（无法删除他人评论）；采用软删除(通过评论id，填充delete_time)

-返回基础响应

  

2.评论列表：GET/douyin/comment/list/(需要调用userRPC)

>query(token,video_id)

- token判断是否过期

- 通过video_id找到delete_time为空(未删除)的所有评论再以评论id倒序返回comments(预加载评论user)

- 返回响应(comments)


  
### relation:
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

- 通过当前用户ID(relation.user_id)和"is_follow=true"查找关系表，找到当前用户关注的所有用户ID

- 并通过userRPC表找出用户信息follows

- 返回响应follows

  

  

3.粉丝列表：GET/douyin/relation/follower/list/

>query(token,user_id)

- token判断是否过期

- 通过当前用户ID(relation.to_user_id)和"is_follow=true"查找关系表，找到关注了当前用户的所有用户ID
- 并通过连接userRPC表找出用户信息fans
- 再查找关注表找到当前用户是否关注了他的粉丝(如果互相关注设定fans的"is_follow"为true)

- 返回响应fans

  

  

4.好友列表：GET/douyin/relation/friend/list/
>query(token,user_id)
- token判断是否过期
- 找到当前用户（user_id）相互关注的用户(user)ID
- 并通过连接userRPC表找出用户信息
- 返回响应friends
  
 ###  message:
 只能和好友（互关）操作
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


#### Redis :

  

#### 旁路缓存模式

  

用户信息(登录操作)

视频信息(上传，feed)

> feed 需要进行排序，用有序集合类型（但随着视频集合元素增多，需要考虑大key问题）

  
  

#### 异步缓存写入：

用户：收到的赞数量

视频：点赞数量

