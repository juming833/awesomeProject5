package model

import (
	"context"
	"fmt"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Conn *gorm.DB

var Rdb *redis.Client

func NewRdb() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.20.16:6379",
		Password: "",
		DB:       0,
	})
	Rdb = rdb
	Store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)
	return
}
func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "123456", "127.0.0.1:3306", "gorm")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库链接错误")
		panic(err)
	}
	// 设置数据表
	Conn = conn
}

var DB *mongo.Client

func NewMongoDB() {
	// 创建MongoDB客户端
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.30.56:27017"))
	if err != nil {
		fmt.Println("无法连接到MongoDB:", err)
		return
	}
	DB = client
}

func Close() {
	db, _ := Conn.DB()

	_ = db.Close()
	_ = Rdb.Close()

}
