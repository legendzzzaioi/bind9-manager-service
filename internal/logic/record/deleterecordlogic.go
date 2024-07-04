package record

import (
	"context"
	"net/http"

	"bind9-manager-service/internal/middleware"
	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRecordLogic {
	return &DeleteRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteRecordLogic) DeleteRecord(r *http.Request, id int) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if id == 0 {
		return &types.Message{Code: 400, Context: "record id cannot be empty"}, nil
	}

	record, err := model.GetRecordById(l.svcCtx.DataSource, id)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	err = model.DeleteRecord(l.svcCtx.DataSource, id)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateOperationLog(l.svcCtx.DataSource, claims.Username, "delete", "record "+record.Domain+" "+record.Name+" "+record.Type+" "+record.Value)

	bindPath := l.svcCtx.Config.BindPath
	if err := model.GenerateZoneFileByDomain(l.svcCtx.DataSource, bindPath, record.Domain); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	if err := svc.ReloadBind9(); err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	return &types.Message{Code: 0, Context: "success"}, nil
}
