package log

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginLogLogic {
	return &GetLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginLogLogic) GetLoginLog(page, pageSize int) (resp *types.LoginLogResp, err error) {
	logs, err := model.GetLoginLog(l.svcCtx.DataSource, page, pageSize)
	if err != nil {
		return nil, err
	}

	count, err := model.GetLoginLogCount(l.svcCtx.DataSource)
	if err != nil {
		return nil, err
	}

	resp = &types.LoginLogResp{
		Logs:  logs,
		Total: count,
	}

	return resp, nil
}
