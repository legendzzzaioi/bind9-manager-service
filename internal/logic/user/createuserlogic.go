package user

import (
	"context"
	"fmt"
	"net/http"

	"bind9-manager-service/internal/middleware"
	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(r *http.Request, req *types.CreateUserReq) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if req.Username == "" || req.Password == "" || req.Role == "" {
		return nil, fmt.Errorf("username or password or role is empty")
	}

	// 创建用户
	err = model.CreateUser(l.svcCtx.DataSource, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateUserLog(l.svcCtx.DataSource, claims.Username, "create", req.Username+" "+req.Role)

	return &types.Message{Code: 0, Context: "success"}, nil
}
