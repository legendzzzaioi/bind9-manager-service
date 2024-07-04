package user

import (
	"context"
	"net/http"

	"bind9-manager-service/internal/middleware"
	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(r *http.Request, username string) (resp *types.Message, err error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	if claims.Role != "admin" {
		return &types.Message{Code: 400, Context: "role forbidden"}, nil
	}

	if username == "" {
		return &types.Message{Code: 400, Context: "username cannot be empty"}, nil
	}

	err = model.DeleteUser(l.svcCtx.DataSource, username)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	model.CreateUserLog(l.svcCtx.DataSource, claims.Username, "delete", username)

	return &types.Message{Code: 0, Context: "success"}, nil
}
