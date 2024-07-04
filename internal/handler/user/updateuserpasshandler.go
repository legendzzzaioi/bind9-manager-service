package user

import (
	"net/http"

	"bind9-manager-service/internal/logic/user"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateUserPassHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserPassReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateUserPassLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserPass(r, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
