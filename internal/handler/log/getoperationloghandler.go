package log

import (
	"net/http"
	"strconv"

	"bind9-manager-service/internal/logic/log"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOperationLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

		l := log.NewGetOperationLogLogic(r.Context(), svcCtx)
		resp, err := l.GetOperationLog(page, pageSize)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
