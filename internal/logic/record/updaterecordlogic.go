package record

import (
	"context"
	"net/http"

	"bind9-manager-service/internal/middleware"
	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRecordLogic {
	return &UpdateRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRecordLogic) UpdateRecord(r *http.Request, req *types.Record) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if req.Id == 0 || req.Domain == "" || req.Name == "" || req.Type == "" || req.Value == "" {
		return &types.Message{Code: 400, Context: "id,domain,name,type,value cannot be empty"}, nil
	}
	err = model.UpdateRecord(l.svcCtx.DataSource, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateOperationLog(l.svcCtx.DataSource, claims.Username, "update", "record "+req.Domain+" "+req.Name+" "+req.Type+" "+req.Value)

	bindPath := l.svcCtx.Config.BindPath
	if err := model.GenerateZoneFileByDomain(l.svcCtx.DataSource, bindPath, req.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
