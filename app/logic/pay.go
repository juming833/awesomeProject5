package logic

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"net/url"
	"os"

	//"github.com/uber/jaeger-client-go/crossdock/client"
	"go-code/awesomeProject1/app/model"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type AlipayConfig struct {
	AppID           string
	AlipayPublicKey string
	PrivateKey      string
}
type Book struct {
	Id    int64
	Title string
	Price float64
}

var client *alipay.Client

const alipayPublicKeyPath = "./file/应用公钥RSA2048.txt"

func HandlePayment(c *gin.Context) {

	config := &AlipayConfig{
		AppID:           "9021000133605979",
		AlipayPublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmQgoaBw/CPUbI4+FNSEGW5E6V5vBU2EZxjLo3Vw+GQmFq9fa/RSVKqcNqm/94jGpRjzIGGvaI+vkKEGFqxn0zR1c8skUCzLVwvCXLKM7bwBjVcxPssz25AqjPOff4naSsGe/KW85D9U33chRr01JTw863k4VQMAED3WKtKOfdmXKNrK31Uc/v41F4Vf9MBQ4YHBFxBMJJlaCVLY9j3MCrFBk/c/DkOqgMhWCRPesZOra07OuLRgvuW8SUKEQ0jGAOXh/wEOPAShsLCVz7kH5SaQvRM0Ajn8UmmvLJeLzCuicmqpBicEGNDZY/UUCk1jkGW88wY4ZUW8ZYj4bFmR9wQIDAQAB",
		PrivateKey:      "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCZCChoHD8I9Rsjj4U1IQZbkTpXm8FTYRnGMujdXD4ZCYWr19r9FJUqpw2qb/3iMalGPMgYa9oj6+QoQYWrGfTNHVzyyRQLMtXC8JcsoztvAGNVzE+yzPbkCqM859/idpKwZ78pbzkP1TfdyFGvTUlPDzreThVAwAQPdYq0o592Zco2srfVRz+/jUXhV/0wFDhgcEXEEwkmVoJUtj2PcwKsUGT9z8OQ6qAyFYJE96xk6trTs64tGC+5bxJQoRDSMYA5eH/AQ48BKGwsJXPuQflJpC9EzQCOfxSaa8sl4vMK6JyaqkGJwQY0Nlj9RQKTWOQZbzzBjhlRbxliPhsWZH3BAgMBAAECggEBAII1icQD/LKPc3cPgCVTo0g4Su8OTTG31ilpWfA3s3K0QnOUJ/XdFIsFN0CFJ/3ViZWE1DDLXAHijfCex7BtVnlMebDu3Momh/HlqovFnCV3rA3t2QZnFc+1KW4CUeIsR+YncGU+CMxq/K0N7Er/tNe7Ori6saS1fP72YX4QFveyk0YIU1mpf4XFHZPJcBtRznIX7RL2tHlUcUEWFuYKNVBhbbEkHHaEPNhkcZvNkQlRw0uB53y38sB/8A8K7WgSWDn6K4Fm8Wo4mRYpE45W10/YIoQqSJrn9TCCmznfn0PRkAbl7bx3ZbeQZwvmt06XRpmzG+a8FcCMs7/lc1kIqUECgYEA35gdbnqx8uiaqYw+Sjq6l2f5JKzkcDXNDcdULRcYpVHjz9EjWAtQZeuSur4UzjzopmA9jd8wl2zZ5phwueUQfOW1kGsuehTJYK1/kT1IiIjzJJJUu2OanyJh1IiDOrfuGINz/7g1zLg+BGXA2+sbkurds/2rHpyS6bBJp0A8KykCgYEArzYChQfMBbeC28FtGzMURWy7UQYyZYpKLBZ26Vd2rfA9IZsx2btSSrUzH5jGpgY+aBVJytQwmeY4xcnsYjS+Lf0SCtaQwEfq4Yo8DJcUgn5Dr5cXQdohzlmwj+ttqWy+RoQWw/ewTKo/+1NpdEcq8ZF3UVOiuTOCtrWixoGpqNkCgYBLIgxRwXzmSO8fpfaQbuegVNnNTyN76Tg8Nwy8EhTACKJqhLZBsZuUg/JMOe831hO/RHuqnqLvXy5hZDwM20INFgtKuhlmAHIiwGorXdJvvsUlNt6uESw16RHCnxSpuICTUCi8dyPs5MRCKgcCXyBm3/EbuFO195sOV6SIRPdy8QKBgEuqGWB1w4fcfiiE2adbYa+xixssf+sV035CnUubt+bZzqyKpvaOTLVylwdYJMFiBOVR/DIZfHuZzn6r+udWO6MEnRwNGSnQZQDNf604OWU5PeXAbAzGNL82QGeGli8KhQlS1bn/ZyVUiJjypqHOch8sYMjAQ4+TTSm2Ovta13fRAoGBAMQS7uQpEzigHrKxSpb1Fua1arF2xLPPnW2JzUuqCi76ahqeSp8liCwKL1Zu8znozi+dpOgbc6DzkSSAVgzD+iSehAINgPDUj9290WG5RuBbGUJKezX/mnk7b9yfUj7NPzsaV2668+HewmZvBhCRmHli83rDleDT5sbmMn2OQeW3",
	}
	// 创建支付宝客户端
	client, err := alipay.New(config.AppID, config.PrivateKey, false)
	// 从全局中间件中获取支付宝客户端
	c.Set("alipay", client)
	client, _ = c.MustGet("alipay").(*alipay.Client)
	var id int64
	idStr := c.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetBook(id)
	// 创建支付请求参数
	p := alipay.TradePagePay{}
	p.NotifyURL = "https://www.baidu.com"    // 设置支付宝回调通知URL
	p.ReturnURL = "https://www.bilibili.com" // 设置支付成功后跳转的URL
	p.Subject = ret.Title
	p.OutTradeNo = strconv.FormatInt(time.Now().Unix(), 10)
	p.TotalAmount = strconv.FormatFloat(ret.Price, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	model.CreateOrder(p.OutTradeNo)
	// 发起支付请求
	if !model.IsOrderOpen(p.OutTradeNo) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单已关闭"})
		return
	}
	result, err := client.TradePagePay(p)

	if err != nil {
		// 处理支付请求错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, result.String())

	go func() {
		// 假设订单超时时间为30分钟
		timeout := time.NewTimer(30 * time.Second)

		<-timeout.C
		// 关闭订单的逻辑
		CloseOrder(p.OutTradeNo)
		// 更新订单状态为关闭状态
		model.UpdateOrderStatus(p.OutTradeNo, "CLOSED")
	}()

}
func CloseOrder(orderNo string) {
	client, _ := alipay.New("9021000133605979", "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCZCChoHD8I9Rsjj4U1IQZbkTpXm8FTYRnGMujdXD4ZCYWr19r9FJUqpw2qb/3iMalGPMgYa9oj6+QoQYWrGfTNHVzyyRQLMtXC8JcsoztvAGNVzE+yzPbkCqM859/idpKwZ78pbzkP1TfdyFGvTUlPDzreThVAwAQPdYq0o592Zco2srfVRz+/jUXhV/0wFDhgcEXEEwkmVoJUtj2PcwKsUGT9z8OQ6qAyFYJE96xk6trTs64tGC+5bxJQoRDSMYA5eH/AQ48BKGwsJXPuQflJpC9EzQCOfxSaa8sl4vMK6JyaqkGJwQY0Nlj9RQKTWOQZbzzBjhlRbxliPhsWZH3BAgMBAAECggEBAII1icQD/LKPc3cPgCVTo0g4Su8OTTG31ilpWfA3s3K0QnOUJ/XdFIsFN0CFJ/3ViZWE1DDLXAHijfCex7BtVnlMebDu3Momh/HlqovFnCV3rA3t2QZnFc+1KW4CUeIsR+YncGU+CMxq/K0N7Er/tNe7Ori6saS1fP72YX4QFveyk0YIU1mpf4XFHZPJcBtRznIX7RL2tHlUcUEWFuYKNVBhbbEkHHaEPNhkcZvNkQlRw0uB53y38sB/8A8K7WgSWDn6K4Fm8Wo4mRYpE45W10/YIoQqSJrn9TCCmznfn0PRkAbl7bx3ZbeQZwvmt06XRpmzG+a8FcCMs7/lc1kIqUECgYEA35gdbnqx8uiaqYw+Sjq6l2f5JKzkcDXNDcdULRcYpVHjz9EjWAtQZeuSur4UzjzopmA9jd8wl2zZ5phwueUQfOW1kGsuehTJYK1/kT1IiIjzJJJUu2OanyJh1IiDOrfuGINz/7g1zLg+BGXA2+sbkurds/2rHpyS6bBJp0A8KykCgYEArzYChQfMBbeC28FtGzMURWy7UQYyZYpKLBZ26Vd2rfA9IZsx2btSSrUzH5jGpgY+aBVJytQwmeY4xcnsYjS+Lf0SCtaQwEfq4Yo8DJcUgn5Dr5cXQdohzlmwj+ttqWy+RoQWw/ewTKo/+1NpdEcq8ZF3UVOiuTOCtrWixoGpqNkCgYBLIgxRwXzmSO8fpfaQbuegVNnNTyN76Tg8Nwy8EhTACKJqhLZBsZuUg/JMOe831hO/RHuqnqLvXy5hZDwM20INFgtKuhlmAHIiwGorXdJvvsUlNt6uESw16RHCnxSpuICTUCi8dyPs5MRCKgcCXyBm3/EbuFO195sOV6SIRPdy8QKBgEuqGWB1w4fcfiiE2adbYa+xixssf+sV035CnUubt+bZzqyKpvaOTLVylwdYJMFiBOVR/DIZfHuZzn6r+udWO6MEnRwNGSnQZQDNf604OWU5PeXAbAzGNL82QGeGli8KhQlS1bn/ZyVUiJjypqHOch8sYMjAQ4+TTSm2Ovta13fRAoGBAMQS7uQpEzigHrKxSpb1Fua1arF2xLPPnW2JzUuqCi76ahqeSp8liCwKL1Zu8znozi+dpOgbc6DzkSSAVgzD+iSehAINgPDUj9290WG5RuBbGUJKezX/mnk7b9yfUj7NPzsaV2668+HewmZvBhCRmHli83rDleDT5sbmMn2OQeW3", false)
	//创建交易关闭请求参数
	closeReq := alipay.TradeClose{
		OutTradeNo: orderNo, // 要查询的订单号
	}
	fmt.Println("Closing order:", orderNo)

	// 发起交易关闭请求
	closeRes, err := client.TradeClose(closeReq)
	if err != nil {
		// 处理交易关闭请求错误
		return
	}

	// 处理交易关闭结果
	if closeRes.Code != "10000" {
		// 交易关闭失败
		return
	}

	// 检查订单状态，如果订单已支付，则不关闭订单
	//if queryRes.TradeStatus == "TRADE_SUCCESS" || queryRes.TradeStatus == "TRADE_FINISHED" {
	//	return
	//}

}

func HandleCallback(c *gin.Context) {
	// 获取请求中的所有参数
	params := make(map[string]string)
	fmt.Println(params)
	c.Request.ParseForm()
	for key, values := range c.Request.Form {
		params[key] = values[0]
	}
	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmQgoaBw/CPUbI4+FNSEGW5E6V5vBU2EZxjLo3Vw+GQmFq9fa/RSVKqcNqm/94jGpRjzIGGvaI+vkKEGFqxn0zR1c8skUCzLVwvCXLKM7bwBjVcxPssz25AqjPOff4naSsGe/KW85D9U33chRr01JTw863k4VQMAED3WKtKOfdmXKNrK31Uc/v41F4Vf9MBQ4YHBFxBMJJlaCVLY9j3MCrFBk/c/DkOqgMhWCRPesZOra07OuLRgvuW8SUKEQ0jGAOXh/wEOPAShsLCVz7kH5SaQvRM0Ajn8UmmvLJeLzCuicmqpBicEGNDZY/UUCk1jkGW88wY4ZUW8ZYj4bFmR9wQIDAQAB"
	// 验证签名
	if VerifySign(params, publicKey) {
		// 签名验证通过，处理业务逻辑
		// TODO: 在这里写下你的业务逻辑代码

		// 返回成功响应
		c.String(http.StatusOK, "success")
	} else {
		// 签名验证失败，返回错误响应
		c.String(http.StatusOK, "error")
	}
}

func VerifySign(params map[string]string, publicKey string) bool {

	// 将参数按照键名进行升序排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接排序后的参数键值对
	var signStr string
	for _, key := range keys {
		if key == "sign" || key == "sign_type" {
			continue
		}
		value := params[key]
		signStr += key + "=" + value + "&"
	}
	signStr = strings.TrimRight(signStr, "&")
	sign := params["sign"]
	fmt.Println(sign)
	// TODO: 进行签名验证的代码
	// 这里需要使用你自己的验签方法，示例中的 Verify 方法仅供参考
	valid := Verify(signStr, sign, publicKey)

	return valid
}

func Verify(signStr, sign, publicKey string) bool {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		// 公钥解码失败
		return false
	}

	// 解析公钥
	pubKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	if err != nil {
		// 公钥解析失败
		return false
	}

	// 计算待签名数据的哈希值
	hashed := sha256.Sum256([]byte(signStr))

	// 解码签名
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		// 签名解码失败
		return false
	}

	// 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		// 签名验证失败
		return false
	}

	// 签名验证通过
	return true
}

func HandleRefund(c *gin.Context) {
	client, _ := alipay.New("9021000133605979", "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCZCChoHD8I9Rsjj4U1IQZbkTpXm8FTYRnGMujdXD4ZCYWr19r9FJUqpw2qb/3iMalGPMgYa9oj6+QoQYWrGfTNHVzyyRQLMtXC8JcsoztvAGNVzE+yzPbkCqM859/idpKwZ78pbzkP1TfdyFGvTUlPDzreThVAwAQPdYq0o592Zco2srfVRz+/jUXhV/0wFDhgcEXEEwkmVoJUtj2PcwKsUGT9z8OQ6qAyFYJE96xk6trTs64tGC+5bxJQoRDSMYA5eH/AQ48BKGwsJXPuQflJpC9EzQCOfxSaa8sl4vMK6JyaqkGJwQY0Nlj9RQKTWOQZbzzBjhlRbxliPhsWZH3BAgMBAAECggEBAII1icQD/LKPc3cPgCVTo0g4Su8OTTG31ilpWfA3s3K0QnOUJ/XdFIsFN0CFJ/3ViZWE1DDLXAHijfCex7BtVnlMebDu3Momh/HlqovFnCV3rA3t2QZnFc+1KW4CUeIsR+YncGU+CMxq/K0N7Er/tNe7Ori6saS1fP72YX4QFveyk0YIU1mpf4XFHZPJcBtRznIX7RL2tHlUcUEWFuYKNVBhbbEkHHaEPNhkcZvNkQlRw0uB53y38sB/8A8K7WgSWDn6K4Fm8Wo4mRYpE45W10/YIoQqSJrn9TCCmznfn0PRkAbl7bx3ZbeQZwvmt06XRpmzG+a8FcCMs7/lc1kIqUECgYEA35gdbnqx8uiaqYw+Sjq6l2f5JKzkcDXNDcdULRcYpVHjz9EjWAtQZeuSur4UzjzopmA9jd8wl2zZ5phwueUQfOW1kGsuehTJYK1/kT1IiIjzJJJUu2OanyJh1IiDOrfuGINz/7g1zLg+BGXA2+sbkurds/2rHpyS6bBJp0A8KykCgYEArzYChQfMBbeC28FtGzMURWy7UQYyZYpKLBZ26Vd2rfA9IZsx2btSSrUzH5jGpgY+aBVJytQwmeY4xcnsYjS+Lf0SCtaQwEfq4Yo8DJcUgn5Dr5cXQdohzlmwj+ttqWy+RoQWw/ewTKo/+1NpdEcq8ZF3UVOiuTOCtrWixoGpqNkCgYBLIgxRwXzmSO8fpfaQbuegVNnNTyN76Tg8Nwy8EhTACKJqhLZBsZuUg/JMOe831hO/RHuqnqLvXy5hZDwM20INFgtKuhlmAHIiwGorXdJvvsUlNt6uESw16RHCnxSpuICTUCi8dyPs5MRCKgcCXyBm3/EbuFO195sOV6SIRPdy8QKBgEuqGWB1w4fcfiiE2adbYa+xixssf+sV035CnUubt+bZzqyKpvaOTLVylwdYJMFiBOVR/DIZfHuZzn6r+udWO6MEnRwNGSnQZQDNf604OWU5PeXAbAzGNL82QGeGli8KhQlS1bn/ZyVUiJjypqHOch8sYMjAQ4+TTSm2Ovta13fRAoGBAMQS7uQpEzigHrKxSpb1Fua1arF2xLPPnW2JzUuqCi76ahqeSp8liCwKL1Zu8znozi+dpOgbc6DzkSSAVgzD+iSehAINgPDUj9290WG5RuBbGUJKezX/mnk7b9yfUj7NPzsaV2668+HewmZvBhCRmHli83rDleDT5sbmMn2OQeW3", false)
	publicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmQgoaBw/CPUbI4+FNSEGW5E6V5vBU2EZxjLo3Vw+GQmFq9fa/RSVKqcNqm/94jGpRjzIGGvaI+vkKEGFqxn0zR1c8skUCzLVwvCXLKM7bwBjVcxPssz25AqjPOff4naSsGe/KW85D9U33chRr01JTw863k4VQMAED3WKtKOfdmXKNrK31Uc/v41F4Vf9MBQ4YHBFxBMJJlaCVLY9j3MCrFBk/c/DkOqgMhWCRPesZOra07OuLRgvuW8SUKEQ0jGAOXh/wEOPAShsLCVz7kH5SaQvRM0Ajn8UmmvLJeLzCuicmqpBicEGNDZY/UUCk1jkGW88wY4ZUW8ZYj4bFmR9wQIDAQAB"
	//err := LoadAliPayPublicKey(publicKey)
	client.LoadAliPayPublicKey(publicKey)
	if err := verifyAlipaySignature(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}
	orderNo := c.Query("order_no") // Get the order number from the request parameters

	refundReq := alipay.TradeRefund{
		OutTradeNo:   orderNo,     // The order number to refund
		RefundAmount: "10.00",     // The refund amount
		OutRequestNo: "refund001", // The unique refund request number
	}

	refundRes, err := client.TradeRefund(refundReq)
	if err != nil {
		// Handle trade refund request error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if refundRes.Code != "10000" {
		// Trade refund failed
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refund failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refund processed successfully"})
}
func verifyAlipaySignature(c *gin.Context) error {
	// 从请求参数中获取待验证的签名字符串和其他参数
	sign := c.Query("sign")
	signType := c.Query("sign_type")
	params := make(map[string]string)
	c.Request.ParseForm()
	for key, values := range c.Request.Form {
		params[key] = values[0]
	}

	// 加载支付宝公钥
	alipayPublicKey, err := loadAlipayPublicKey(alipayPublicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load Alipay public key: %v", err)
	}

	// 提取签名字符串并进行 URL 解码
	signStr := extractSignString(params)
	signStr, err = urlDecode(signStr)
	if err != nil {
		return fmt.Errorf("failed to URL decode sign string: %v", err)
	}

	// 对签名字符串进行验签
	if err := verifySignature(signStr, sign, signType, alipayPublicKey); err != nil {
		return fmt.Errorf("failed to verify Alipay signature: %v", err)
	}

	return nil
}

// 加载支付宝公钥
func loadAlipayPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	// 读取公钥文件内容
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Alipay public key file: %v", err)
	}

	// 解析 PEM 格式的公钥
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing Alipay public key")
	}

	// 解析公钥
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Alipay public key: %v", err)
	}

	// 转换为 RSA 公钥类型
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("alipay public key is not an RSA public key")
	}

	return rsaPublicKey, nil
}

// 提取签名字符串
func extractSignString(params map[string]string) string {
	var keys []string
	for key := range params {
		if key != "sign" && key != "sign_type" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var signItems []string
	for _, key := range keys {
		value := params[key]
		signItems = append(signItems, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(signItems, "&")
}

// URL 解码
func urlDecode(str string) (string, error) {
	decoded, err := url.QueryUnescape(str)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

// 验证签名
func verifySignature(signStr, sign, signType string, publicKey *rsa.PublicKey) error {
	// 根据签名类型选择哈希算法
	var hash crypto.Hash
	switch signType {
	case "RSA2":
		hash = crypto.SHA256
	case "RSA":
		hash = crypto.SHA1
	default:
		return fmt.Errorf("unsupported sign_type: %s", signType)
	}

	// 对签名进行 Base64 解码
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	// 计算签名摘要
	h := hash.New()
	if _, err := h.Write([]byte(signStr)); err != nil {
		return fmt.Errorf("failed to compute hash: %v", err)
	}
	hashed := h.Sum(nil)

	// 验证签名
	if err := rsa.VerifyPKCS1v15(publicKey, hash, hashed, signature); err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	return nil
}