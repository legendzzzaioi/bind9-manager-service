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

func (l *DeleteZoneLogic) DeleteZone(r *http.Request, domain string, record bool) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if domain == "" {
		return &types.Message{Code: 400, Context: "domain cannot be empty"}, nil
	}
	err = model.DeleteZone(l.svcCtx.DataSource, domain, record)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateOperationLog(l.svcCtx.DataSource, claims.Username, "delete", "zone "+domain)

	bindPath := l.svcCtx.Config.BindPath

	if err := model.GenerateNamedLocalConf(l.svcCtx.DataSource, bindPath); err != nil {
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
