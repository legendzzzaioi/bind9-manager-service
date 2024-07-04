package log

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogLogic {
	return &GetUserLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogLogic) GetUserLog(page, pageSize int) (resp *types.UserLogResp, err error) {
	logs, err := model.GetUserLog(l.svcCtx.DataSource, page, pageSize)
	if err != nil {
		return nil, err
	}

	count, err := model.GetUserLogCount(l.svcCtx.DataSource)
	if err != nil {
		return nil, err
	}

	resp = &types.UserLogResp{
		Logs:  logs,
		Total: count,
	}

	return resp, nil
}
