// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserInfoResponse struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"username"`
}
