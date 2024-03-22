package handler

import (
	"go-code/awesomeProject1/app/login_zero/common/response"
	"net/http"

	"go-code/awesomeProject1/app/login_zero/api/internal/logic"
	"go-code/awesomeProject1/app/login_zero/api/internal/svc"
)

func userInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo()
		response.Response(r, w, resp, err)
	}
}
