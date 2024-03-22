package model

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBorrow2(t *testing.T) {
	var wg sync.WaitGroup
	concurrency := 3 // 并发事务的数量
	successCount := 0

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		userId := int64(i + 1) // 分别使用 userId 1, 2, 3
		go func() {
			defer wg.Done()

			// 在每个事务中尝试同时更新同一行数据
			err := Borrow2(userId, 1) // 使用 Borrow 函数作为示例
			if err == nil {
				successCount++
			}
		}()
	}

	wg.Wait()

	// 检查只有一个事务成功更新数据
	if successCount != 1 {
		t.Errorf("Expected 1 successful transaction, got %d", successCount)
	}
}
func TestNewMongoDB(t *testing.T) {
	NewMongoDB()
	// 检查全局变量DB是否为空
	if DB == nil {
		fmt.Println("MongoDB连接未初始化")
		return
	}
	// 等待连接完成
	time.Sleep(2 * time.Second)
	// 检查连接状态
	err := DB.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("MongoDB连接失败:", err)
		return
	}
	fmt.Println("成功连接到MongoDB")
}
func TestStoreImagesInMongoDB(t *testing.T) {
	NewMongoDB()
	imageDir := "E:\\go.code\\src\\go-code\\awesomeProject3\\app\\images"
	collectionName := "book"
	StoreImagesInMongoDB(imageDir, collectionName)
}
