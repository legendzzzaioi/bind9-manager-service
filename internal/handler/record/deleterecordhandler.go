package record

import (
	"errors"
	"net/http"
	"strconv"

	"bind9-manager-service/internal/logic/record"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			httpx.ErrorCtx(r.Context(), w, errors.New("id cannot be empty"))
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.New("id must be an integer"))
			return
		}

		l := record.NewDeleteRecordLogic(r.Context(), svcCtx)
		resp, err := l.DeleteRecord(id)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
