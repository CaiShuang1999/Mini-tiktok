package video

import (
	"net/http"

	"tiktok-micro/app/api/internal/logic/video"
	"tiktok-micro/app/api/internal/svc"
	"tiktok-micro/app/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublishRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewPublishLogic(r.Context(), svcCtx)
		resp, err := l.Publish(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
