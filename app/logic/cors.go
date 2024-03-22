package logic

import "github.com/gin-gonic/gin"

// CorsMiddleware CORS中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT,POST,DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,x-requested-with")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
func OptionsHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT,POST,DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,x-requested-with")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")

	c.AbortWithStatus(200)
}

//
//func HandleProxy(c *gin.Context) {
//	targetURL := "https://openapi-sandbox.dl.alipaydev.com" + c.Param("path")
//	remote, err := url.Parse(targetURL)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// 创建反向代理
//	proxy := httputil.NewSingleHostReverseProxy(remote)
//	proxy.ModifyResponse = func(resp *http.Response) error {
//		resp.Header.Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080") // 设置CORS头部信息
//		return nil
//	}
//
//	// 修改请求头部，以便正确地转发请求
//	c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api")
//	c.Request.URL.Host = remote.Host
//	c.Request.URL.Scheme = remote.Scheme
//	c.Request.Host = remote.Host
//
//	// 转发请求
//	proxy.ServeHTTP(c.Writer, c.Request)
//}
