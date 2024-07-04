package login

import (
	"net/http"

	"bind9-manager-service/internal/logic/login"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		l := login.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout(r, ip)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
