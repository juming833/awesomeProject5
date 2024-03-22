package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-code/awesomeProject1/app/login_zero/api/internal/svc"
	"go-code/awesomeProject1/app/login_zero/api/internal/types"
	"go-code/awesomeProject1/app/login_zero/common/jwts"
	"go-code/awesomeProject1/app/login_zero/rpc/types/user"
)

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

func (l *LoginLogic) Login(req *types.LoginRequest) (string, error) {
	loginResp, err := l.svcCtx.UserRpc.UserLogin(l.ctx, &user.UserLoginRequest{
		Username: req.UserName,
		Password: req.Password,
	})
	fmt.Println(loginResp)
	user, _ := l.svcCtx.UsersModel.FindOneByUsername(l.ctx, req.UserName)
	if user == nil {
		return "", errors.New("用户名或密码错误")
	}

	// 比较密码
	if user.Password != req.Password {
		return "", errors.New("用户名或密码错误")
	}
	// 生成 JWT token
	auth := l.svcCtx.Config.Auth
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		Username: req.UserName,
	}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		return "", err
	}
	return token, nil
}
