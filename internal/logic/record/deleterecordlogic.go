package record

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRecordLogic {
	return &DeleteRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRecordLogic) DeleteRecord(id int) (resp *types.Message, err error) {
	if id == 0 {
		return &types.Message{Code: 400, Context: "record id cannot be empty"}, nil
	}

	record, err := model.GetRecordById(l.svcCtx.DB, id)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	err = model.DeleteRecord(l.svcCtx.DB, id)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	bindPath := l.svcCtx.Config.BindPath
	if err := model.GenerateZoneFileByDomain(l.svcCtx.DB, bindPath, record.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
