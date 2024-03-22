package svc

import (
	"go-code/awesomeProject1/app/login_zero/common/init_gorm"
	"go-code/awesomeProject1/app/login_zero/common/models"
	"go-code/awesomeProject1/app/login_zero/rpc/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := init_db.InitGorm(c.Mysql.DataSource)
	db.AutoMigrate(&models.UserModel{})
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
