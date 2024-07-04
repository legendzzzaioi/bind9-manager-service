package user

import (
	"net/http"

	"bind9-manager-service/internal/logic/user"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")

		l := user.NewDeleteUserLogic(r.Context(), svcCtx)
		resp, err := l.DeleteUser(r, username)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
