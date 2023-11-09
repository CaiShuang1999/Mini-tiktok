package feed

import (
	"net/http"

	"tiktok-micro/app/api/internal/logic/feed"
	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FeedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FeedRequest

		if err := httpx.Parse(r, &req); err != nil {

			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := feed.NewFeedLogic(r.Context(), svcCtx)
		resp, err := l.Feed(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
