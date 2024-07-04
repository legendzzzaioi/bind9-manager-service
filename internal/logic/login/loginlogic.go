package login

import (
	"context"
	"fmt"
	"time"

	"bind9-manager-service/internal/model"
	"bind9-manager-service/internal/svc"
	"bind9-manager-service/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type MyClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(ip string, req *types.LoginReq) (resp *types.LoginResp, err error) {
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}

	// 验证用户名密码
	password, err := model.GetPasswordByName(l.svcCtx.DataSource, req.Username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// 写入登陆日志
	model.CreateLoginLog(l.svcCtx.DataSource, req.Username, ip, "login")

	// 获取user信息
	user, err := model.GetUserByName(l.svcCtx.DataSource, req.Username)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	expire := time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire) * time.Second
	secret := []byte(l.svcCtx.Config.JwtAuth.AccessSecret)

	// 填充claims
	claims := MyClaims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Unix(now, 0)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}

	// 生成jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	resp = &types.LoginResp{
		Token: tokenString,
	}

	return
}
