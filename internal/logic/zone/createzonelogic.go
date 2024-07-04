package zone

import (
	"context"
	"net/http"

	"bind9-manager-service/internal/middleware"
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

func (l *CreateZoneLogic) CreateZone(r *http.Request, req *types.ZoneReq) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if req.Domain == "" || req.Ttl == 0 || req.CacheTtl == 0 || req.Expire == 0 || req.MailAddress == "" || req.PrimaryNameServer == "" || req.Refresh == 0 || req.Retry == 0 {
		return &types.Message{Code: 400, Context: "domain,ttl,cache_ttl,expires,mail_address,primary_name_server,refresh,retry,ttl cannot be empty"}, nil
	}
	err = model.CreateZone(l.svcCtx.DataSource, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateOperationLog(l.svcCtx.DataSource, claims.Username, "create", "zone "+req.Domain)

	bindPath := l.svcCtx.Config.BindPath

	if err := model.GenerateNamedLocalConf(l.svcCtx.DataSource, bindPath); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := model.GenerateZoneFileByDomain(l.svcCtx.DataSource, bindPath, req.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
