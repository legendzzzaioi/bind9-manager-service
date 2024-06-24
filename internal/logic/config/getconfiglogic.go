package config

import (
	"context"
	"fmt"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConfigLogic {
	return &GetConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConfigLogic) GetConfig(key string) (resp *types.Config, err error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}
	config, err := model.GetConfig(l.svcCtx.DB, key)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
