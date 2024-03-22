package logic

//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/smartwalle/alipay/v3"
//	"log"
//	"net/http"
//	"time"
//)
//
//var aliClient *alipay.Client
//
//const (
//	kServerPort   = "8080"
//	kServerDomain = "http://xx.xx.xx.xx" + ":" + kServerPort
//)
//
//func main() {
//	var err error
//
//	if aliClient, err = alipay.New(cert.AppId, cert.PrivateKey, false); err != nil {
//		log.Println("初始化支付宝失败", err)
//		return
//	}
//
//	// 使用支付宝证书
//	if err = aliClient.LoadAppPublicCertFromFile("cert/appCertPublicKey_2021000118628779.crt"); err != nil {
//		log.Println("加载证书发生错误", err)
//		return
//	}
//	if err = aliClient.LoadAliPayRootCertFromFile("cert/alipayRootCert.crt"); err != nil {
//		log.Println("加载证书发生错误", err)
//		return
//	}
//	if err = aliClient.LoadAliPayPublicCertFromFile("cert/alipayCertPublicKey_RSA2.crt"); err != nil {
//		log.Println("加载证书发生错误", err)
//		return
//	}
//
//	var s = gin.Default()
//	s.GET("/alipay", pay)
//	s.GET("/callback", callback)
//	s.POST("/notify", notify)
//	s.Run(":" + kServerPort)
//}
//
//func pay(c *gin.Context) {
//	var tradeNo = fmt.Sprintf("%d", xid.Next())
//
//	var p = alipay.TradePagePay{}
//	p.Subject = "支付宝测试:" + tradeNo
//	p.OutTradeNo = time.Now().Format("20060102") + "_" + "name"
//	p.TotalAmount = "0.1"
//	url, _ := aliClient.TradePagePay(p)
//
//	c.Redirect(http.StatusTemporaryRedirect, url.String())
//}
//
//func callback(c *gin.Context) {
//	c.Request.ParseForm()
//
//	ok, err := aliClient.VerifySign(c.Request.Form)
//	if err != nil {
//		log.Println("回调验证签名发生错误", err)
//		return
//	}
//
//	if ok == false {
//		log.Println("回调验证签名未通过")
//		return
//	}
//
//	fmt.Println("------c.Request.Form: -------\n", c.Request.Form)
//	var outTradeNo = c.Request.Form.Get("out_trade_no")
//	var p = alipay.TradeQuery{}
//	p.OutTradeNo = outTradeNo
//	rsp, err := aliClient.TradeQuery(p)
//	if err != nil {
//		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s", outTradeNo, err.Error())
//		return
//	}
//	if rsp.IsSuccess() == false {
//		c.String(http.StatusBadRequest, "验证订单 %s 信息发生错误: %s-%s", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
//		return
//	}
//
//	c.String(http.StatusOK, "订单 %s 支付成功", outTradeNo)
//}
//
//func notify(c *gin.Context) {
//	fmt.Println("------notify------")
//	c.Request.ParseForm()
//
//	fmt.Println("c.Request.Form: ", c.Request.Form)
//	ok, err := aliClient.VerifySign(c.Request.Form)
//	if err != nil {
//		log.Println("异步通知验证签名发生错误", err)
//		return
//	}
//
//	if ok == false {
//		log.Println("异步通知验证签名未通过")
//		return
//	}
//
//	log.Println("异步通知验证签名通过")
//
//	var outTradeNo = c.Request.Form.Get("out_trade_no")
//	var p = alipay.TradeQuery{}
//	p.OutTradeNo = outTradeNo
//	rsp, err := aliClient.TradeQuery(p)
//	if err != nil {
//		log.Printf("异步通知验证订单 %s 信息发生错误: %s \n", outTradeNo, err.Error())
//		return
//	}
//	if rsp.IsSuccess() == false {
//		log.Printf("异步通知验证订单 %s 信息发生错误: %s-%s \n", outTradeNo, rsp.Content.Msg, rsp.Content.SubMsg)
//		return
//	}
//
//	log.Printf("订单 %s 支付成功 \n", outTradeNo)
//}
