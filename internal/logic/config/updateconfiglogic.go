package config

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConfigLogic {
	return &UpdateConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConfigLogic) UpdateConfig(req *types.Config) (resp *types.Message, err error) {
	if req.Key == "" || req.Value == "" {
		return &types.Message{Code: 400, Context: "key,value cannot be empty"}, nil
	}
	err = model.UpdateConfig(l.svcCtx.DB, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	err = model.GenerateNamedOptionsConf(l.svcCtx.DB, l.svcCtx.Config.BindPath)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
