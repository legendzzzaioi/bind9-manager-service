package record

import (
	"context"
	"fmt"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecordsLogic {
	return &GetRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecordsLogic) GetRecords(domain string) (resp []types.Record, err error) {
	if domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}
	records, err := model.GetRecords(l.svcCtx.DB, domain)
	if err != nil {
		return nil, err
	}
	return records, nil
}
