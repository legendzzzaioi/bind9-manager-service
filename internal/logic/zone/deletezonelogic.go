package zone

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteZoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteZoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteZoneLogic {
	return &DeleteZoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteZoneLogic) DeleteZone(domain string, record bool) (resp *types.Message, err error) {
	if domain == "" {
		return &types.Message{Code: 400, Context: "domain cannot be empty"}, nil
	}
	err = model.DeleteZone(l.svcCtx.DB, domain, record)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	bindPath := l.svcCtx.Config.BindPath

	if err := model.GenerateNamedLocalConf(l.svcCtx.DB, bindPath); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := model.DeleteZoneFileByDomain(bindPath, domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
