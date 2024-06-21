package zone

import (
	"bind9-manager-service/internal/model"
	"context"

	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetZonesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetZonesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetZonesLogic {
	return &GetZonesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetZonesLogic) GetZones() (resp []types.Zone, err error) {
	zones, err := model.GetZones(l.svcCtx.DB)
	if err != nil {
		return nil, err
	}
	return zones, nil
}
