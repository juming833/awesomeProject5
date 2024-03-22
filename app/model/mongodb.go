package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"path/filepath"
)

func StoreImagesInMongoDB(imageDir string, collectionName string) {
	// 获取集合
	collection := DB.Database("1").Collection(collectionName)
	// 读取目录下的所有图片文件
	files, err := os.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}
	// 遍历文件列表并将图片存储到MongoDB中
	for _, file := range files {
		filePath := filepath.Join(imageDir, file.Name())

		// 读取图片文件内容
		imageData, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file: %s, error: %s\n", filePath, err)
			continue
		}

		// 创建文档并插入到集合中
		document := bson.M{
			"filename": file.Name(),
			"image":    imageData,
		}
		_, err = collection.InsertOne(context.TODO(), document)
		if err != nil {
			log.Printf("Failed to insert document: %s, error: %s\n", file.Name(), err)
		} else {
			fmt.Printf("Inserted document: %s\n", file.Name())
		}
	}
}
