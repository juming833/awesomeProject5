package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-code/awesomeProject1/app/logic"
	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/tools"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func New(port int) {
	r := gin.Default()
	r.LoadHTMLGlob("app/view/*")
	r.Static("/img", "app/images")
	r.Use(logic.CorsMiddleware())
	r.GET("/index", logic.Index)
	r.GET("/books", logic.GetBooks)
	r.GET("/cover", logic.Cover)
	r.POST("/email", logic.Email)
	r.GET("/sms", logic.SendSms)
	r.POST("/verifyCode", logic.VerifyCode)
	r.POST("/verifyCode2", logic.VerifyCodeHandler)
	i := r.Group("")
	x := r.Group("")

	//i.Use(checkStudent)
	//x.Use(checkAdmin)
	{
		//login
		r.GET("/login", logic.GetLogin)
		r.GET("/logout", logic.Logout)
		r.POST("/login", logic.PostLogin)
		//r.GET("/wechat", logic.CheckSignature)
		//r.GET("/wechat/login", logic.Redirect)
		//r.GET("/wechat/callback", logic.Callback)

		r.GET("/adminLogin", logic.GetAdminLogin)
		r.POST("/adminLogin", logic.AdminPostLogin)
		r.GET("/create", logic.GetCreate)
		r.POST("/createUser", logic.CreateUser)
		r.GET("/email-login", logic.EmailLogin)
		r.GET("/phone-login", logic.PhoneLogin)
		r.POST("/borrow", logic.Borrow)
		r.POST("/return", logic.ReturnBook)
		r.GET("/record", logic.GetRecord)

	}
	{
		//restful
		r.POST("/book", logic.AddBook)
		r.DELETE("/book", logic.DelBook)
		r.PUT("/book", logic.UpdateBook)

	}
	{
		i.GET("/bookInfo", logic.GetBookInfo)
		i.POST("/book/add", logic.AddBook)       //1
		i.POST("/book/update", logic.UpdateBook) //2
		i.POST("/book/del", logic.DelBook)
		i.GET("/records", logic.GetRecords)
	}
	{
		x.POST("/addAdmin", logic.AddAdmin)
		x.POST("/updateAdmin", logic.UpdateAdmin)
		x.POST("/delAdmin", logic.DelAdmin)
	}
	{

		//r.Any("/api/*path", logic.HandleProxy)

		r.GET("/pay", logic.HandlePayment)
		r.OPTIONS("/pay", logic.OptionsHandler)
		r.POST("/callback", logic.HandleCallback)
		r.GET("/refund", logic.HandleRefund)
	}
	{
		r.MaxMultipartMemory = 8 << 20 // 8 MiB
		r.POST("/upload", logic.Upload)
		r.POST("/uploads", logic.FilesUpload)
	}
	//验证码
	{
		r.GET("/captcha", logic.GetCaptcha)

		r.POST("/captcha/verify", func(context *gin.Context) {
			var param tools.CaptchaData
			if err := context.ShouldBind(&param); err != nil {
				context.JSON(http.StatusOK, tools.ParamErr)
				return
			}

			fmt.Printf("参数为：%+v", param)
			if !tools.CaptchaVerify(param) {
				context.JSON(http.StatusOK, tools.ECode{
					Code:    10008,
					Message: "验证失败",
				})
				return
			}
			context.JSON(http.StatusOK, tools.OK)
		})
	}
	if err := r.Run(":8080"); err != nil {
		panic("gin 启动失败")
	}

	//gin优雅退出

	//server := &http.Server{
	//	Addr:    ":8080",
	//	Handler: r,
	//}
	//
	//go func() {
	//	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		panic("gin 启动失败: " + err.Error())
	//	}
	//}()
	//gracefulShutdown(server)
	//// 无限循环保持程序运行
	//for {
	//	time.Sleep(time.Second)
	//}

}
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("sever started")
	<-quit
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
func checkStudent(context *gin.Context) {
	var name string
	var id int64
	var roleId int64
	values := model.GetSession(context)
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id, _ = strconv.ParseInt(v.(string), 10, 64)
	}
	if v, ok := values["role_id"]; ok {
		roleId = v.(int64)
	}
	fmt.Println(roleId)
	if name == "" || id <= 0 || roleId != 2 {
		context.JSON(http.StatusUnauthorized, tools.ECode{
			Message: "You do not have permission to access this route.",
		})
		context.Abort()

	}

	context.Next()
}
func checkAdmin(context *gin.Context) {
	var name string
	var id int64
	var roleId int64
	values := model.GetSession(context)
	if v, ok := values["name"]; ok {
		name = v.(string)
	}
	if v, ok := values["id"]; ok {
		id, _ = strconv.ParseInt(v.(string), 10, 64)
	}
	if v, ok := values["role_id"]; ok {
		roleId = v.(int64)
	}
	fmt.Println(roleId)
	if name == "" || id <= 0 || roleId != 3 {
		context.JSON(http.StatusUnauthorized, tools.ECode{
			Message: "You do not have permission to access this route.",
		})
		context.Abort()

	}

	context.Next()
}
