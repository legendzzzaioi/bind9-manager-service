package zone

import (
	"context"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateZoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateZoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateZoneLogic {
	return &CreateZoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateZoneLogic) CreateZone(req *types.ZoneReq) (resp *types.Message, err error) {
	if req.Domain == "" || req.Ttl == 0 || req.CacheTtl == 0 || req.Expire == 0 || req.MailAddress == "" || req.PrimaryNameServer == "" || req.Refresh == 0 || req.Retry == 0 {
		return &types.Message{Code: 400, Context: "domain,ttl,cache_ttl,expires,mail_address,primary_name_server,refresh,retry,ttl cannot be empty"}, nil
	}
	err = model.CreateZone(l.svcCtx.DB, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	bindPath := l.svcCtx.Config.BindPath

	if err := model.GenerateNamedLocalConf(l.svcCtx.DB, bindPath); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := model.GenerateZoneFileByDomain(l.svcCtx.DB, bindPath, req.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
