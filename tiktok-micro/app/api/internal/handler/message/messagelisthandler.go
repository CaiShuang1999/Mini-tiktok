package message

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tiktok-micro/app/api/internal/logic/message"
	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"
)

func MessageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := message.NewMessageListLogic(r.Context(), svcCtx)
		resp, err := l.MessageList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
