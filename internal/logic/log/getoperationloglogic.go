package log

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperationLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperationLogLogic {
	return &GetOperationLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOperationLogLogic) GetOperationLog(page, pageSize int) (resp *types.OperationLogResp, err error) {
	logs, err := model.GetOperationLog(l.svcCtx.DataSource, page, pageSize)
	if err != nil {
		return nil, err
	}

	count, err := model.GetOperationLogCount(l.svcCtx.DataSource)
	if err != nil {
		return nil, err
	}

	resp = &types.OperationLogResp{
		Logs:  logs,
		Total: count,
	}

	return resp, nil
}
