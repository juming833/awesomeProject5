package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"

	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/tools"
	"net/http"
)

//用户                          系统
//|                             |
//|---------发送登录请求--------->|
//|                             |
//|<-------返回登录页面----------|
//|                             |
//|---------提交登录凭证--------->|
//|                             |
//|<-------验证登录凭证----------|
//|                             |
//|--------返回登录成功---------->|
//|                             |

type Admin struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key" json:"id"`
	RoleId      int64     `gorm:"column:role_id;type:bigint(20)" json:"role_id"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Password    string    `gorm:"column:password;type:varchar(255)" json:"password" form:"password"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *Admin) TableName() string {
	return "admin"
}

type User struct {
	Id           int64  `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Uid          int64  `gorm:"column:uid;type:bigint(20)" json:"uid"`
	Name         string `json:"name" form:"name"`
	Password     string `json:"password" form:"password"`
	Email        string `gorm:"column:email;type:varchar(255)" json:"email"`
	Phone        string `gorm:"column:phone;type:varchar(255)" json:"phone"`
	CaptchaId    string `json:"captcha_id"form:"captcha_id"`
	CaptchaValue string `json:"captcha_value"form:"captcha_value"`
}

func GetLogin(context *gin.Context) {
	context.HTML(200, "login.html", nil)
}
func GetAdminLogin(context *gin.Context) {
	context.HTML(200, "admin.html", nil)
}
func EmailLogin(context *gin.Context) {
	context.HTML(200, "email.html", nil)
}
func PhoneLogin(context *gin.Context) {
	context.HTML(200, "phone.html", nil)
}
func AdminPostLogin(context *gin.Context) {
	var admin Admin
	if err := context.ShouldBind(&admin); err != nil {
		context.JSON(200, tools.ECode{
			Message: err.Error(),
		})
	}
	ret := model.GetAdmin(admin.Name)
	if ret.Id < 1 || ret.Password != admin.Password {
		context.JSON(200, tools.UserErr)
		return
	}
	_ = model.SetSession(context, admin.Name, ret.Id, ret.RoleId)
	context.JSON(200, tools.ECode{
		Message: "登陆成功",
	})
	return
}
func PostLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBind(&user); err != nil {
		context.JSON(200, tools.ECode{
			Message: err.Error(),
		})
	}
	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败！",
		})
		return
	}
	ret := model.GetStudent(user.Name)
	if ret.Id < 1 || ret.Password != EncryptV1(user.Password) {
		context.JSON(200, tools.UserErr)
		return
	}
	jwt, _ := model.GetJwt(ret.Id, user.Name)
	context.SetCookie("jwt", jwt, 3600, "/", "", true, false)
	_ = model.SetSession(context, user.Name, ret.Id, ret.RoleId)
	context.JSON(200, tools.ECode{
		Message: "登陆成功",
	})
	return

}

func Logout(context *gin.Context) {
	_ = model.FlushSession(context)
	context.Redirect(302, "/login")
}

type CStudent struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
	Phone     string `json:"phone"`
}

type CAdmin struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
}

func GetCreate(context *gin.Context) {
	context.HTML(200, "create.html", nil)
}

func CreateUser(context *gin.Context) {
	var student CStudent
	if err := context.ShouldBind(&student); err != nil {
		context.JSON(200, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}

	if student.Name == "" || student.Password == "" || student.Password2 == "" {
		context.JSON(http.StatusOK, tools.ParamErr)
		return
	}

	//校验密码
	if student.Password != student.Password2 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同！", //这里有风险
		})
		return
	}

	nameLen := len(student.Name)
	password := len(student.Password)
	phone := len(student.Phone)
	if nameLen > 16 || nameLen < 6 || password > 16 || password < 6 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "账号或密码大于6小于16",
		})
		return
	}
	if phone != 11 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "手机号格式错误",
		})
		return
	}

	//密码不能是纯数字 -》 数字+小写字母+大写字母
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(student.Password) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字", //这里有风险
		})
		return
	}

	//这里有一个巨大的BUG，并发安全！
	if oldStudent := model.GetStudent(student.Name); oldStudent.Id > 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在！",
		})
		return
	}

	Student := model.User{
		Name:        student.Name,
		Password:    EncryptV1(student.Password),
		Phone:       EncryptV2(student.Phone),
		UpdatedTime: time.Now(),
		CreatedTime: time.Now(),
	}

	if err := model.CreateUser(&Student); err != nil {
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

func EncryptV1(pwd string) string {
	// 创建一个 MD5 哈希对象
	hash := md5.New()
	// 将密码转换为字节数组并计算哈希值
	hash.Write([]byte(pwd))
	hashBytes := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)
	// 打印加密后的密码
	fmt.Printf("加密后的密码：%s\n", hashString)
	// 返回加密后的密码字符串
	return hashString
}
func EncryptV2(phone string) string {
	if len(phone) != 11 {
		// 如果手机号长度不等于11位，则不进行脱敏处理
		return phone
	}
	desensitized := phone[:3] + "****" + phone[7:]
	return desensitized
}
