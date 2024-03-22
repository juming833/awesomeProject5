package logic

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"go-code/awesomeProject1/app/model"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// SendSms 发送短信验证码
func SendSms(c *gin.Context) {
	// 从请求中获取手机号码
	phoneNumber := c.Query("phone_number")

	// 生成随机验证码
	code := strconv.Itoa(rand.Intn(900000) + 100000)
	model.Rdb.Set(c, "verification_code", code, 60*time.Second)

	// 阿里云短信API的相关配置信息
	accessKeyId := "LTAI5t6GvbKJCGPjZV58Nq2p"
	accessSecret := "6jXChqRCYy9rXFXO3J2W6fZZhp4vHu"
	signName := "阿里云短信测试"
	templateCode := "SMS_154950909"

	// 创建短信客户端
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建发送短信请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumber
	request.SignName = signName
	request.TemplateCode = templateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	// 发送短信
	response, err := client.SendSms(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 判断短信发送是否成功
	if response.Code != "OK" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": response.Message})
		return
	}
	// 返回成功信息
	c.JSON(http.StatusOK, gin.H{"message": "短信发送成功"})
}
func VerifyCodeHandler(c *gin.Context) {
	type requestData struct {
		VerificationCode string `form:"verificationCode" binding:"required"`
	}

	var req requestData
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数有误"})
		return
	}

	//phoneNumber := req.PhoneNumber
	inputVerificationCode := req.VerificationCode
	verificationCode, err := model.Rdb.Get(c, "verification_code").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证码已过期或不存在"})
		return
	}

	if inputVerificationCode == verificationCode {
		// 验证通过，返回登录成功
		c.JSON(http.StatusOK, gin.H{"success": true})
	} else {
		// 验证失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码输入错误"})
	}
}
