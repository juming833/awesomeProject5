package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-code/awesomeProject1/app/login_zero/api/internal/svc"
)

type RenderLoginHtmlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenderLoginHtmlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenderLoginHtmlLogic {
	return &RenderLoginHtmlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenderLoginHtmlLogic) RenderLoginHtml() (resp string, err error) {
	// todo: add your logic here and delete this line

	return
}
