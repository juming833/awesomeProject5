package response

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type Body struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func Response(r *http.Request, w http.ResponseWriter, res any, err error) {
	if err != nil {
		body := Body{
			Code: 10086,
			Data: nil,
			Msg:  err.Error(),
		}
		httpx.WriteJson(w, 200, body)
		return
	}
	body := Body{
		Code: 0,
		Data: res,
		Msg:  "成功",
	}
	httpx.WriteJson(w, 200, body)

}
