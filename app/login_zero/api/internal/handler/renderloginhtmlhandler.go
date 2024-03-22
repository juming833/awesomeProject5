package handler

import (
	"html/template"
	"net/http"
)

func renderLoginHandler(w http.ResponseWriter, r *http.Request) {
	// 读取 login.html 文件

	tmpl, err := template.ParseFiles("E:\\go.code\\src\\go-code\\awesomeProject5\\app\\view\\login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 渲染并发送响应
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
