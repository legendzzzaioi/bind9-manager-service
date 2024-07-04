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

type UpdateUserPassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserPassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPassLogic {
	return &UpdateUserPassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserPassLogic) UpdateUserPass(r *http.Request, req *types.UpdateUserPassReq) (resp *types.Message, err error) {
	// 从上下文中获取claims
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*middleware.MyClaims)
	if !ok {
		return &types.Message{Code: 400, Context: "unauthorized"}, nil
	}

	// 检查请求中的用户名和密码是否为空
	if req.Username == "" || req.Password == "" {
		return &types.Message{Code: 400, Context: "username or password is empty"}, nil
	}

	// 权限检查：普通用户只能修改自己的密码，管理员可以修改所有用户的密码
	if claims.Username != req.Username && claims.Role != "admin" {
		return &types.Message{Code: 403, Context: "forbidden: only admin can update other user's password"}, nil
	}

	// 更新用户密码
	err = model.UpdateUserPass(l.svcCtx.DataSource, *req)
	if err != nil {
		return &types.Message{Code: 400, Context: err.Error()}, nil
	}

	// 记录用户操作日志
	model.CreateUserLog(l.svcCtx.DataSource, claims.Username, "update", req.Username+" password")

	// 返回成功响应
	return &types.Message{Code: 0, Context: "success"}, nil
}
