package logic

import (
	"context"

	"go-code/awesomeProject1/app/login_zero/rpc/internal/svc"
	"go-code/awesomeProject1/app/login_zero/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLoginLogic {
	return &UserLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserLoginLogic) UserLogin(in *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	// todo: add your logic here and delete this line

	return &user.UserLoginResponse{}, nil
}
