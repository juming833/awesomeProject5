package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rbcervilla/redisstore/v9"
)

var Store *redisstore.RedisStore
var SessionName = "session-name"

func GetSession(c *gin.Context) map[interface{}]interface{} {
	session, _ := Store.Get(c.Request, SessionName)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSession(c *gin.Context, name string, id int64, roleId int64) error {
	session, _ := Store.Get(c.Request, SessionName)
	session.Values["name"] = name
	session.Values["id"] = id
	session.Values["role_id"] = roleId
	return session.Save(c.Request, c.Writer)
}

func FlushSession(c *gin.Context) error {
	session, _ := Store.Get(c.Request, SessionName)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = ""
	session.Values["role_id"] = ""
	session.Options.MaxAge = -1
	c.SetCookie("name", "", -1, "/", "", true, false)
	c.SetCookie("Id", "", -1, "/", "", true, false)
	c.SetCookie("jwt", "", -1, "/", "", true, false)
	return session.Save(c.Request, c.Writer)
}
