package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-code/awesomeProject1/app/model"
	"go-code/awesomeProject1/app/tools"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddBook(context *gin.Context) {
	idStr := context.Query("title")
	author, _ := context.GetPostForm("author")
	count, _ := context.GetPostForm("count")
	num, _ := strconv.ParseInt(count, 10, 32)
	//构建结构体
	Book := model.BookInfo{
		Title:       idStr,
		Author:      author,
		Count:       int(num),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := model.AddBook(Book); err != nil {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, tools.OK)
	return
}
func UpdateBook(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	count, _ := context.GetPostForm("count")
	num, _ := strconv.ParseInt(count, 10, 32)
	book := model.BookInfo{
		Id:          id,
		Count:       int(num),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	if err := model.UpdateBook(book); err != nil {
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

func DelBook(context *gin.Context) {
	var id int64
	idStr := context.Query("id")
	id, _ = strconv.ParseInt(idStr, 10, 64)
	if err := model.DelBook(id); err != true {
		context.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "删除失败",
		})
		return
	}
	context.JSON(200, tools.OK)
	return
}

func Cover(c *gin.Context) {
	filename := c.Query("filename")
	// 假设你有一个函数用于根据文件名加载图书封面图片的字节数据
	coverData, err := model.LoadCoverImage(filename)
	if err != nil {
		// 处理加载图片数据失败的情况
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	// 设置响应的 Content-Type 为 image/jpeg（或其他适当的 MIME 类型）
	c.Header("Content-Type", "image/jpeg")

	// 将图书封面图片数据作为响应的内容返回给客户端
	c.Data(http.StatusOK, "image/jpeg", coverData)
}
func Upload(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	dst := "./file/" + file.Filename
	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(file, dst)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
func FilesUpload(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload[]"] // 注意这里名字不要对不上了
	for _, file := range files {
		log.Println(file.Filename)
		// 上传文件至指定目录
		c.SaveUploadedFile(file, "./file/"+file.Filename)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}
