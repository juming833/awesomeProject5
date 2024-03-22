package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go-code/awesomeProject1/app/login_zero/api/internal/config"
	"go-code/awesomeProject1/app/login_zero/model"
	"go-code/awesomeProject1/app/login_zero/rpc/users"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    users.Users
	UsersModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:     c,
		UserRpc:    users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		UsersModel: model.NewUserModel(mysqlConn),
	}
}
