package logic

import (
	"github.com/gin-gonic/gin"
	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/tools"
	"net/http"
	"strconv"
	"time"
)

func AddAdmin(context *gin.Context) {
	var admin CAdmin
	if err := context.ShouldBind(&admin); err != nil {
		context.JSON(200, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}

	if admin.Name == "" || admin.Password == "" || admin.Password2 == "" {
		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}

	//校验密码
	if admin.Password != admin.Password2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同！", //这里有风险
		})
		return
	}

	nameLen := len(admin.Name)
	password := len(admin.Password)
	if nameLen > 16 || nameLen < 6 || password > 16 || password < 6 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "账号或密码大于6小于16",
		})
		return
	}
	//这里有一个巨大的BUG，并发安全！
	if oldAdmin := model.GetAdmin(admin.Name); oldAdmin.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在！",
		})
		return
	}

	Admin := model.Admin{
		Name:        admin.Name,
		Password:    admin.Password,
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	if err := model.AddAdmin(&Admin); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10007,
			Message: "新用户创建失败！", //这里有风险
		})
		return
	}

	context.JSON(http.StatusOK, tools.ECode{

		Message: "创建成功",
	})

	return
}

func UpdateAdmin(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	password, _ := context.GetPostForm("password")
	admin := model.Admin{
		Id:          id,
		Password:    password,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := model.UpdateAdmin(admin); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": 10006,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
	return
}

func DelAdmin(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	if err := model.DelAdmin(id); err != true {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "删除失败",
		})
		return
	}
	context.JSON(200, tools.OK)
	return
}
