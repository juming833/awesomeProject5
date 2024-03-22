package logic

//
//import (
//	"context"
//	"crypto/sha1"
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"github.com/redis/go-redis/v9"
//	"github.com/skip2/go-qrcode"
//	"go.uber.org/zap"
//	"io"
//	"net/http"
//	"net/url"
//	"sort"
//	"time"
//)
//
//// TOKEN 假设您在Go代码中定义了一个名为TOKEN的常量，用于存储您的令牌值
//const TOKEN = "111"
//
//var Client redis.Client
//
//// 配置公众号的token
//
//func CheckSignature(c *gin.Context) {
//	signature := c.Query("signature")
//	timestamp := c.Query("timestamp")
//	nonce := c.Query("nonce")
//	echostr := c.Query("echostr")
//
//	tmpArr := []string{TOKEN, timestamp, nonce}
//	sort.Strings(tmpArr)
//	tmpStr := ""
//	for _, v := range tmpArr {
//		tmpStr += v
//	}
//
//	tmpHash := sha1.New()
//	tmpHash.Write([]byte(tmpStr))
//	tmpStr = fmt.Sprintf("%x", tmpHash.Sum(nil))
//	fmt.Println(tmpStr)
//	fmt.Println(signature)
//
//	if tmpStr == signature {
//		c.String(http.StatusOK, echostr)
//		ctx := context.Background()
//		redis.Client.Set(Client, ctx, "library:token", tmpStr, 7*24*time.Hour)
//	} else {
//		c.String(http.StatusForbidden, "Signature verification failed "+timestamp)
//	}
//}
//
//// Redirect 微信扫码登录
//// @Summary 用户登录接口3
//// @Description 通过微信扫码登录，手机进行登录验证
//// @Tags 公开
//// @Accept json
//// @Produce application/json
//// @Param Url query string true "内网穿透地址"
//// @Router /api/v1/wechat/login [get]
//func Redirect(c *gin.Context) {
//	path := c.Query("Url") //hr37hu.natappfree.cc
//
//	state := uuid.New().String()[:5]                                     //防止跨站请求伪造攻击 增加安全性
//	redirectURL := url.QueryEscape("http://" + path + "wechat/callback") //userinfo,
//	fmt.Println(redirectURL)
//	wechatLoginURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=%s&scope=snsapi_userinfo#wechat_redirect", "wx2cb4b8d66dea007b", redirectURL, state)
//	wechatLoginURL, _ = url.QueryUnescape(wechatLoginURL)
//	// 生成二维码
//	qrCode, err := qrcode.Encode(wechatLoginURL, qrcode.Medium, 256)
//	if err != nil {
//		// 错误处理
//		c.String(http.StatusInternalServerError, "Error generating QR code")
//		return
//	}
//	// 将二维码图片作为响应返回给用户
//	c.Header("Content-Type", "image/png")
//	c.Writer.Write(qrCode)
//}
//
//type ResponseData struct {
//	Data    interface{}
//	Message string
//	Code    interface{}
//}
//
//func Callback(c *gin.Context) {
//	// 获取微信返回的授权码
//	code := c.Query("code")
//	// 向微信服务器发送请求，获取access_token和openid
//	tokenResp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", "wx2cb4b8d66dea007b", "c6e95614a454d965f0f8ccd0d9177c3c", code))
//	if err != nil {
//		fmt.Println(err)
//		resp := &ResponseData{
//			Data:    nil,
//			Message: "error,获取token失败",
//			Code:    "",
//		}
//		c.JSON(http.StatusBadRequest, resp)
//		return
//	}
//	// 解析响应中的access_token和openid
//	var tokenData struct {
//		AccessToken  string `json:"access_token"`
//		ExpiresIn    int    `json:"expires_in"`
//		RefreshToken string `json:"refresh_token"`
//		OpenID       string `json:"openid"`
//		Scope        string `json:"scope"`
//	}
//	if err1 := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err1 != nil {
//		resp := &ResponseData{
//			Data:    nil,
//			Message: "error,获取token失败",
//			Code:    "",
//		}
//		c.JSON(http.StatusBadRequest, resp)
//		return
//	}
//	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", tokenData.AccessToken, tokenData.OpenID)
//	userInfoResp, err := http.Get(userInfoURL)
//	if err != nil {
//		// 错误处理
//		zap.L().Error("获取失败")
//		return
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			fmt.Println("未知错误")
//		}
//	}(userInfoResp.Body)
//	zap.L().Info("登录成功")
//	return
//}
