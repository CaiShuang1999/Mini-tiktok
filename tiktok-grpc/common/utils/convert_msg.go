package utils

import (
	"tiktok-grpc/apps/comment/comment_pb"
	"tiktok-grpc/apps/favorite/favorite_pb"
	"tiktok-grpc/apps/feed/feed_pb"
	"tiktok-grpc/apps/message/messages_pb"
	"tiktok-grpc/apps/relation/relation_pb"
	"tiktok-grpc/apps/user/pb"
	"tiktok-grpc/apps/video/video_pb"
	"tiktok-grpc/model"
)

func ConvertUserToProto(user model.User) *pb.UserInfo {
	protoUser := &pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}

func ConvertAuthorToProto(user model.User) *video_pb.UserInfo {
	protoUser := &video_pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}

func ConvertVideoToProto(video model.Video) *video_pb.VideoInfo {
	protoVideo := &video_pb.VideoInfo{
		Id:            video.ID,
		Author:        ConvertAuthorToProto(video.Author),
		CommentCount:  video.CommentCount,
		CoverUrl:      video.CoverURL,
		PlayUrl:       video.PlayURL,
		Title:         video.Title,
		CreateTime:    video.CreateTime,
		FavoriteCount: video.FavoriteCount,
		IsFavorite:    video.IsFavorite,
	}

	return protoVideo
}

// /feed
func ConvertAuthorTofeedProto(user model.User) *feed_pb.UserInfo {
	protoUser := &feed_pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}

func ConvertVideoTofeedProto(video model.Video) *feed_pb.VideoInfo {
	protoVideo := &feed_pb.VideoInfo{
		Id:            video.ID,
		Author:        ConvertAuthorTofeedProto(video.Author),
		CommentCount:  video.CommentCount,
		CoverUrl:      video.CoverURL,
		PlayUrl:       video.PlayURL,
		Title:         video.Title,
		CreateTime:    video.CreateTime,
		FavoriteCount: video.FavoriteCount,
		IsFavorite:    video.IsFavorite,
	}

	return protoVideo
}

// favorite

func ConvertAuthorTofavoriteProto(user model.User) *favorite_pb.UserInfo {
	protoUser := &favorite_pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}

func ConvertVideoTofavoriteProto(video model.Video) *favorite_pb.VideoInfo {
	protoVideo := &favorite_pb.VideoInfo{
		Id:            video.ID,
		Author:        ConvertAuthorTofavoriteProto(video.Author),
		CommentCount:  video.CommentCount,
		CoverUrl:      video.CoverURL,
		PlayUrl:       video.PlayURL,
		Title:         video.Title,
		CreateTime:    video.CreateTime,
		FavoriteCount: video.FavoriteCount,
		IsFavorite:    video.IsFavorite,
	}

	return protoVideo
}

// 评论

func ConvertUserToCommentProto(user model.User) *comment_pb.UserInfo {
	protoUser := &comment_pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}

func ConvertCommentToProto(comment model.Comment) *comment_pb.Comment {
	protoComment := &comment_pb.Comment{
		Id:         comment.ID,
		User:       ConvertUserToCommentProto(comment.UserMsg),
		Content:    comment.CommentMsg,
		CreateDate: comment.CreateDate,
	}

	return protoComment
}

// 关系

func ConvertUserToRelationProto(user model.User) *relation_pb.UserInfo {
	protoUser := &relation_pb.UserInfo{

		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}

	return protoUser
}

// 聊天记录

func ConvertMessageToMessageProto(message model.Message) *messages_pb.MessageInfo {
	protoMessageInfo := &messages_pb.MessageInfo{
		Id:         message.ID,
		ToUserId:   message.ToUserID,
		FromUserId: message.FromUserID,
		Content:    message.Content,
		CreateTime: message.CreateTime,
	}

	return protoMessageInfo
}
