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

type CreateRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRecordLogic {
	return &CreateRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRecordLogic) CreateRecord(r *http.Request, req *types.CreateRecord) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if req.Domain == "" || req.Name == "" || req.Type == "" || req.Value == "" {
		return &types.Message{Code: 400, Context: "domain,name,type,value cannot be empty"}, nil
	}
	err = model.CreateRecord(l.svcCtx.DataSource, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateOperationLog(l.svcCtx.DataSource, claims.Username, "create", "record "+req.Domain+" "+req.Name+" "+req.Type+" "+req.Value)

	bindPath := l.svcCtx.Config.BindPath
	if err := model.GenerateZoneFileByDomain(l.svcCtx.DataSource, bindPath, req.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
