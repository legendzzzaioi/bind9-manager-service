package zone

import (
	"net/http"

	"bind9-manager-service/internal/logic/zone"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteZoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		domain := r.URL.Query().Get("domain")
		record := r.URL.Query().Get("record") == "true"
		l := zone.NewDeleteZoneLogic(r.Context(), svcCtx)
		resp, err := l.DeleteZone(domain, record)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
