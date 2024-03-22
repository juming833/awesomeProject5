package app

import (
	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/router"
)

func Start() {
	model.NewMysql()
	model.NewRdb()
	defer func() {
		model.Close()
	}()
	port := 8080
	router.New(port)
}
