package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"go-code/awesomeProject1/app/login_zero/api/internal/svc"
	"go-code/awesomeProject1/app/login_zero/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	// todo: add your logic here and delete this line

	userId := l.ctx.Value("user_id").(json.Number)
	fmt.Printf("%v, %T, \n", userId, userId)
	username := l.ctx.Value("username").(string)
	uid, _ := userId.Int64()
	return &types.UserInfoResponse{
		UserId:   uid,
		UserName: username,
	}, nil
}
