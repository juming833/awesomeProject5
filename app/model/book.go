package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"time"
)

func AddBook(book BookInfo) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&book).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func UpdateBook(book BookInfo) error {
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&BookInfo{}).Where("id = ?", book.Id).Updates(book).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
func DelBook(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&BookInfo{}, id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}

	return true
}

func GetBooks() []BookInfo {
	ret := make([]BookInfo, 0)
	if err := Conn.Table("book_info").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

//func GetBooks(limit, offset int) []BookInfo {
//	ret := make([]BookInfo, 0)
//	if err := Conn.Table("book_info").Limit(limit).Offset(offset).Find(&ret).Error; err != nil {
//		fmt.Printf("err:%s", err.Error())
//	}
//	return ret
//}

func GetRecords() []StudentBook {
	ret := make([]StudentBook, 0)
	if err := Conn.Table("student_book").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func GetBook(id int64) BookInfo {
	var ret BookInfo
	if err := Conn.Table("book_info").Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}
func GetRecord(id int64) []StudentBook {
	var ret []StudentBook
	if err := Conn.Table("student_book").Where("user_id = ?", id).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

func Borrow(userId int64, name string, Id int64) error {
	var count int
	var studentBook StudentBook
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("student_book").Where("user_id = ? AND name= ? AND book_info_id = ?", userId, name, Id).First(&studentBook).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			// 如果存在记录且状态为0，则已经借过，不能再借
			if studentBook.Status == 0 {
				return errors.New("已经借过此书，不能再借")
			}
		}
		if err := db.Table("book_info").Where("id = ? AND count > 0", Id).Updates(map[string]interface{}{
			"count": gorm.Expr("count  - ?", 1),
		}).Error; err != nil {
			return err
		}

		if err := db.Table("book_info").Where("id = ?", Id).Select("count").Scan(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			user := StudentBook{
				UserId:      userId,
				Name:        name,
				BookInfoid:  Id,
				Status:      0,
				CreatedTime: time.Now(),
				UpdateTime:  time.Now(),
			}
			if err := db.Create(&user).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
			}
			return nil
		}
		return errors.New("库存不足")
	})
	return err
}
func ReturnBook(userId int64, Id int64) error {
	var count int
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("book_info").Where("id= ?", Id).Select("count").Scan(&count).Error; err != nil {
			return err
		}
		count++
		if err := db.Table("book_info").Where("id = ?", Id).Update("count", count).Error; err != nil {
			return err
		}
		if err := db.Table("student_book").Where("user_id=? AND book_info_id = ? AND status = 0", userId, Id).Update("status", 1).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

func LoadCoverImage(filename string) ([]byte, error) {
	filePath := "E:\\go.code\\src\\go-code\\awesomeProject1\\app\\images\\" + filename
	fmt.Println(filePath)
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func GetUid(Id int64) (int64, error) {
	user := User{}
	err := Conn.Transaction(func(db *gorm.DB) error {
		if err := db.Table("user").Where("id = ?", Id).Select("uid").Scan(&user).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return user.Uid, nil
}

// Borrow2 悲观锁
func Borrow2(userId int64, Id int64) error {
	var count int
	var studentBook StudentBook
	var ret BookInfo
	err := Conn.Transaction(func(db *gorm.DB) error {
		// 使用悲观锁查询book_info表的记录，并锁定该行数据
		if err := db.Table("book_info").Where("id = ?", Id).Set("count", "FOR UPDATE").First(&ret).Error; err != nil {
			return err
		}

		if err := db.Table("student_book").Where("user_id = ? AND book_info_id = ?", userId, Id).First(&studentBook).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			// 如果存在记录且状态为0，则已经借过，不能再借
			if studentBook.Status == 0 {
				return errors.New("已经借过此书，不能再借")
			}
		}

		if ret.Count >= 1 {
			count = ret.Count - 1

			// 更新book_info表中的count字段
			if err := db.Table("book_info").Where("id = ?", Id).Update("count", count).Error; err != nil {
				return err
			}

			user := StudentBook{
				UserId:      userId,
				BookInfoid:  Id,
				Status:      0,
				CreatedTime: time.Now(),
				UpdateTime:  time.Now(),
			}
			if err := db.Create(&user).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
			}
			return nil
		}
		return errors.New("库存不足")
	})
	return err
}

// Borrow3 乐观锁
func Borrow3(userId int64, Id int64) error {
	var count int
	var studentBook StudentBook
	err := Conn.Transaction(func(db *gorm.DB) error {
		// 查询学生借书记录
		if err := db.Table("student_book").Where("user_id = ? AND book_info_id = ?", userId, Id).First(&studentBook).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			// 如果存在记录且状态为0，则已经借过，不能再借
			if studentBook.Status == 0 {
				return errors.New("已经借过此书，不能再借")
			}
		}

		// 更新图书信息表，同时检查 updateTime 字段以实现乐观锁
		result := db.Table("book_info").
			Where("id = ? AND count > 0 AND update_time = ?", Id, studentBook.UpdateTime).
			Updates(map[string]interface{}{
				"count":      gorm.Expr("count - ?", 1),
				"updateTime": time.Now(),
			})
		if result.Error != nil {
			return result.Error
		}

		// 检查更新的行数，如果为0，则表示数据已被其他事务修改，当前事务需要回滚
		if result.RowsAffected == 0 {
			return errors.New("图书信息已被修改，请刷新重试")
		}

		// 重新查询图书的库存
		if err := db.Table("book_info").Where("id = ?", Id).Select("count").Scan(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			// 创建新的学生借书记录
			user := StudentBook{
				UserId:      userId,
				BookInfoid:  Id,
				Status:      0,
				CreatedTime: time.Now(),
				UpdateTime:  time.Now(),
			}
			if err := db.Create(&user).Error; err != nil {
				return err
			}
			return nil
		}

		return errors.New("库存不足")
	})

	return err
}
