package logic

import (
	"crypto/tls"
	"fmt"
	"go-code/awesomeProject1/app/model"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Email(context *gin.Context) {
	// 邮件服务器的地址和端口
	smtpHost := "smtp.qq.com"
	smtpPort := 587

	// 发件人的认证信息
	senderEmail := "1738353843@qq.com"
	senderPassword := "ugpgqhngcgzrdgig"

	// 收件人的地址z
	recipientEmail, _ := context.GetPostForm("email")
	//var user User
	//if err := context.ShouldBindJSON(&user); err != nil {
	//	context.JSON(200, tools.ECode{
	//		Message: err.Error(),
	//	})
	//}
	//fmt.Println(user)
	// 通过收件人电子邮件获取用户ID
	//userId := model.GetUserByEmail(recipientEmail)
	//fmt.Println(userId)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		context.JSON(http.StatusBadRequest, gin.H{
	//			"error": "用户不存在",
	//		})
	//		return
	//	}
	//	context.JSON(http.StatusInternalServerError, gin.H{
	//		"error": "无法获取用户ID",
	//	})
	//	return
	//}
	ret := model.GetUserByEmail(recipientEmail)
	//if ret.Id < 1 {
	//	context.JSON(200, tools.UserErr)
	//	return
	//}
	fmt.Println(ret)
	fmt.Println(ret.Id, ret.Name)
	JWT, _ := model.GetJwt(ret.Id, ret.Name)
	fmt.Println(JWT)
	context.SetCookie("jwt", JWT, 3600, "/", "", true, false)
	_ = model.SetSession(context, ret.Name, ret.Id, ret.RoleId)
	//model.Rdb.Set(context, "id", userId, 30*time.Second)
	// 生成随机验证码
	verificationCode := strconv.Itoa(rand.Intn(900000) + 100000)
	// 将验证码存储到 Redis 中，有效期为 5 分钟
	model.Rdb.Set(context, "verification_code", verificationCode, 60*time.Second)
	// 邮件内容
	subject := "验证码登录"
	body := "您的验证码是：" + verificationCode

	// 组装邮件体
	message := "From: " + senderEmail + "\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// 配置SMTP客户端
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// 连接到邮件服务器
	smtpAddr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	client, err := smtp.Dial(smtpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	if err = client.StartTLS(&tls.Config{ServerName: smtpHost}); err != nil {
		log.Fatal(err)
	}

	// 发起身份验证
	if err = client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	// 设置发件人
	if err = client.Mail(senderEmail); err != nil {
		log.Fatal(err)
	}

	// 设置收件人
	if err = client.Rcpt(recipientEmail); err != nil {
		log.Fatal(err)
	}

	// 发送邮件内容
	dataWriter, err := client.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer dataWriter.Close()

	_, err = fmt.Fprintf(dataWriter, message)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
	})
}

func VerifyCode(context *gin.Context) {
	var requestData struct {
		VerificationCode string `form:"verificationCode" binding:"required"`
	}

	if err := context.ShouldBind(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "请求参数有误"})
		return
	}

	inputVerificationCode := requestData.VerificationCode
	verificationCode, err := model.Rdb.Get(context, "verification_code").Result()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "验证码已过期或不存在"})
		return
	}

	// 检查验证码是否正确
	if inputVerificationCode == verificationCode {
		// 验证通过，返回登录成功
		context.JSON(http.StatusOK, gin.H{"success": true})
	} else {
		// 验证失败，返回错误信息
		context.JSON(http.StatusBadRequest, gin.H{"error": "验证码输入错误"})
	}
}
