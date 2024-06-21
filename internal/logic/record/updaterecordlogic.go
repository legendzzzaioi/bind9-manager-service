package record

import (
	"context"

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

func (l *UpdateRecordLogic) UpdateRecord(req *types.Record) (resp *types.Message, err error) {
	if req.Id == 0 || req.Domain == "" || req.Name == "" || req.Type == "" || req.Value == "" {
		return &types.Message{Code: 400, Context: "id,domain,name,type,value cannot be empty"}, nil
	}
	err = model.UpdateRecord(l.svcCtx.DB, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	bindPath := l.svcCtx.Config.BindPath
	if err := model.GenerateZoneFileByDomain(l.svcCtx.DB, bindPath, req.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
