package zone

import (
	"net/http"

	"bind9-manager-service/internal/logic/zone"
	"bind9-manager-service/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetZonesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := zone.NewGetZonesLogic(r.Context(), svcCtx)
		resp, err := l.GetZones()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
