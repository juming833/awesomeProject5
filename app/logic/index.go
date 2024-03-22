package logic

import (
	"crypto/md5"
	"encoding/json"
	"strings"

	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/tools"

	"net/http"
	"strconv"
	"time"
)

func Index(context *gin.Context) {
	ret := model.GetBooks()
	context.HTML(http.StatusOK, "index.html", gin.H{"book": ret})
}

func GetBookInfo(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetBook(id)
	context.JSON(200, tools.ECode{
		Data: ret,
	})
}
func GetRecord(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	ret := model.GetRecord(id)
	context.JSON(200, tools.ECode{
		Data: ret,
	})
}

func GetBooks(context *gin.Context) {
	limitStr := context.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	// 获取offset参数，默认为0
	offsetStr := context.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	// 调用数据库或其他数据存储机制获取图书数据
	books := model.GetBooks()

	// 对图书数据进行分页处理
	paginatedBooks := paginateBooks(books, limit, offset)
	jsonData, err := json.Marshal(paginatedBooks)
	if err != nil {
		return
	}
	key := fmt.Sprintf("paginated_data:%d:%d", offset, limit)
	err = model.Rdb.Set(context, key, jsonData, 0).Err()
	if err != nil {
		return
	}
	model.Rdb.Expire(context, key, 30*time.Second)
	// 返回分页结果
	context.JSON(200, tools.ECode{
		Data: paginatedBooks,
	})

}
func paginateBooks(books []model.BookInfo, limit, offset int) []model.BookInfo {
	start := offset
	end := offset + limit
	// 处理边界情况，确保不超过图书列表的长度
	if start >= len(books) {
		return nil
	}
	if end > len(books) {
		end = len(books)
	}

	return books[start:end]
}

func Borrow(context *gin.Context) {
	//if !UidXyz(context) {
	//	context.JSON(http.StatusOK, tools.ECode{
	//		Code:    10016,
	//		Message: "单身二十年",
	//	})
	//	return
	//}
	bookId, _ := context.GetPostForm("bookId")
	Id, _ := strconv.ParseInt(bookId, 10, 64)
	cacheKey := fmt.Sprintf("borrow:%d", Id)
	jwt, _ := context.Cookie("jwt")
	JWT, err := model.CheckJwt(jwt)
	if err != nil {
		context.JSON(200, tools.ECode{
			Message: "检验失败",
		})
		return
	}
	userId := JWT.Id
	name := JWT.Name
	if userId < 0 {
		context.JSON(404, tools.ECode{
			Message: "您还未登录",
		})
		context.Redirect(302, "/login")
		return
	}

	_, err = model.Rdb.Get(context, cacheKey).Result()
	if err == nil {
		context.JSON(200, tools.ECode{
			Message: fmt.Sprintf("你已经借过这本书了"),
		})
		return
	}

	if err := model.Borrow(userId, name, Id); err != nil {
		context.JSON(200, tools.ECode{
			Message: "库存不足",
		})
		return
	}
	cacheValue := fmt.Sprintf("%d-%d", userId, Id)
	model.Rdb.Set(context, cacheKey, cacheValue, 0)
	context.JSON(200, tools.ECode{
		Message: "借阅成功",
	})
	return
}

func ReturnBook(context *gin.Context) {
	//if !UidXyz(context) {
	//	context.JSON(http.StatusOK, tools.ECode{
	//		Code:    10016,
	//		Message: "单身二十年",
	//	})
	//	return
	//}
	IdStr, _ := context.GetPostForm("bookId")
	Id, _ := strconv.ParseInt(IdStr, 10, 64)
	jwt, _ := context.Cookie("jwt")
	JWT, err := model.CheckJwt(jwt)
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "检验失败",
		})
		return
	}
	userIDStr := JWT.Id
	if userIDStr < 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "未登录",
		})
	}
	cacheKey := fmt.Sprintf("borrow:%d", Id)
	cacheResult, err := model.Rdb.Get(context, cacheKey).Result()
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "还书失败，未找到借书记录",
		})
		return
	}

	// 检查借书记录的用户ID是否匹配当前用户
	borrowUserID, _ := strconv.ParseInt(strings.Split(cacheResult, "-")[0], 10, 64)
	if borrowUserID != userIDStr {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "还书失败，借书记录不匹配当前用户",
		})
		return
	}

	// 删除Redis中的借书记录
	if err := model.Rdb.Del(context, cacheKey).Err(); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "还书失败，删除借书记录出错",
		})
		return
	}

	// 执行还书操作
	if err := model.ReturnBook(userIDStr, Id); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "还书失败",
		})
		return
	}

	context.JSON(http.StatusOK, tools.ECode{
		Message: "还书成功",
	})
	return
}

// CheckXYZ 限流
func CheckXYZ(context *gin.Context) bool {
	ip := context.ClientIP()
	ua := context.GetHeader("user-agent")
	fmt.Printf("ip:%s\n,ua:%s\n", ip, ua)

	//转为MD5
	hash := md5.New()           //创建一个MD5哈希实例
	hash.Write([]byte(ip + ua)) //将IP地址和user-agent信息拼接后写入哈希实例。
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes) //将哈希值转换为字符串

	flag, _ := model.Rdb.Get(context, "ban-"+hashString).Bool()
	if flag {
		return false
	}
	i, _ := model.Rdb.Get(context, "xyz-"+hashString).Int() // 从Redis中获取"xyz-"+hashString键对应的值，并将其转换为整数类型
	if i > 5 {
		model.Rdb.SetEx(context, "ban-"+hashString, true, 30*time.Second)
		return false
	}
	// 如果获取的值大于5，则将"ban-"+hashString键设置为true（加入黑名单），并设置过期时间为30秒
	model.Rdb.Incr(context, "xyz-"+hashString)                  //Incr将存储值递增一，Expire用于设置过期时间
	model.Rdb.Expire(context, "xyz-"+hashString, 5*time.Second) //每次访问时次数加一，并设置过期时间5秒
	return true

}
func GetCaptcha(context *gin.Context) {
	if !CheckXYZ(context) {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10016,
			Message: "单身二十年",
		})
		return
	}
	captcha, err := tools.CaptchaGenerate()
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, tools.ECode{
		Data: captcha,
	})
}
func UidXyz(context *gin.Context) bool {
	session, _ := model.Store.Get(context.Request, model.SessionName)
	Id, ok := session.Values["id"].(int64)
	if !ok {
		return ok
	}
	uid, _ := model.GetUid(Id)
	fmt.Println(uid)
	if uid != 0 {
		uidStr := strconv.FormatInt(uid, 10)
		uidHash := md5.New()
		uidHash.Write([]byte(uidStr))
		uidHashBytes := uidHash.Sum(nil)
		uidHashString := hex.EncodeToString(uidHashBytes)
		fmt.Println(uidHashString)
		uidFlag, _ := model.Rdb.Get(context, "ban-"+uidHashString).Bool()
		if uidFlag {
			return false
		}
		uidCount, _ := model.Rdb.Get(context, "xyz-"+uidHashString).Int()
		if uidCount > 5 {
			model.Rdb.SetEx(context, "ban-"+uidHashString, true, 30*time.Second)
			return false
		}

		model.Rdb.Incr(context, "xyz-"+uidHashString)
		model.Rdb.Expire(context, "xyz-"+uidHashString, 5*time.Second)
		return true

	} else {
		ip := context.ClientIP()
		ua := context.GetHeader("user-agent")
		fmt.Printf("ip:%s\n,ua:%s\n", ip, ua)

		//转为MD5
		hash := md5.New()           //创建一个MD5哈希实例
		hash.Write([]byte(ip + ua)) //将IP地址和user-agent信息拼接后写入哈希实例。
		hashBytes := hash.Sum(nil)
		hashString := hex.EncodeToString(hashBytes) //将哈希值转换为字符串

		flag, _ := model.Rdb.Get(context, "ban-"+hashString).Bool()
		if flag {
			return false
		}
		i, _ := model.Rdb.Get(context, "xyz-"+hashString).Int() // 从Redis中获取"xyz-"+hashString键对应的值，并将其转换为整数类型
		if i > 5 {
			model.Rdb.SetEx(context, "ban-"+hashString, true, 30*time.Second)
			return false
		}
		// 如果获取的值大于5，则将"ban-"+hashString键设置为true（加入黑名单），并设置过期时间为30秒
		model.Rdb.Incr(context, "xyz-"+hashString)                  //Incr将存储值递增一，Expire用于设置过期时间
		model.Rdb.Expire(context, "xyz-"+hashString, 5*time.Second) //每次访问时次数加一，并设置过期时间5秒
		return true
	}
}
func GetRecords(context *gin.Context) {
	jwt, _ := context.Cookie("jwt")
	JWT, err := model.CheckJwt(jwt)
	fmt.Println(jwt)
	if err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "检验失败",
		})
		return
	}
	userIDStr := JWT.Id
	if userIDStr < 0 {
		context.JSON(http.StatusOK, tools.ECode{
			Message: "未登录",
		})
	}
	record := model.GetRecords()
	context.JSON(200, tools.ECode{
		Data: record,
	})
}
