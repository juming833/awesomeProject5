package model

import "time"

type BookInfo struct {
	Id                 int64     `gorm:"column:id;type:int(11) unsigned;primary_key;comment:书的id" json:"id"`
	Uid                int64     `gorm:"column:uid;default:NULL"json:"uid"`
	Title              string    `gorm:"column:title;type:varchar(200);comment:书名" json:"title"`
	Author             string    `gorm:"column:author;type:varchar(50);comment:作者" json:"author"`
	Count              int       `gorm:"column:count;type:int(11)" json:"count"`
	PublishingHouse    string    `gorm:"column:publishing_house;type:varchar(50);comment:出版社" json:"publishing_house"`
	Translator         string    `gorm:"column:translator;type:varchar(50);comment:译者" json:"translator"`
	PublishDate        time.Time `gorm:"column:publish_date;type:date;comment:出版时间" json:"publish_date"`
	Pages              int       `gorm:"column:pages;type:int(11);default:100;comment:页数" json:"pages"`
	ISBN               string    `gorm:"column:ISBN;type:varchar(20);comment:ISBN号码" json:"ISBN"`
	Price              float64   `gorm:"column:price;type:double;default:1;comment:价格" json:"price"`
	BriefIntroduction  string    `gorm:"column:brief_introduction;type:varchar(15000);comment:内容简介" json:"brief_introduction"`
	AuthorIntroduction string    `gorm:"column:author_introduction;type:varchar(5000);comment:作者简介" json:"author_introduction"`
	ImgUrl             string    `gorm:"column:img_url;type:varchar(200);comment:封面地址" json:"img_url"`
	DelFlg             int       `gorm:"column:del_flg;type:int(11);default:0;comment:删除标识" json:"del_flg"`
	CreatedTime        time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime        time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *BookInfo) TableName() string {
	return "book_info"
}

type User struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Uid         int64     `gorm:"column:uid;type:bigint(20)" json:"uid"`
	RoleId      int64     `gorm:"column:role_id;type:bigint(20)" json:"role_id"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Password    string    `gorm:"column:password;type:varchar(255)" json:"password"`
	Email       string    `gorm:"column:email;type:varchar(255)" json:"email"`
	Phone       string    `gorm:"column:phone;type:varchar(255)" json:"phone"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *User) TableName() string {
	return "user"
}

type StudentBook struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	UserId      int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	BookInfoid  int64     `gorm:"column:book_info_id;type:bigint(20)" json:"book_info_id"`
	Status      int       `gorm:"column:status;type:int(11)" json:"status"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

func (m *StudentBook) TableName() string {
	return "student_book"
}

type Role struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	Role        string    `gorm:"column:role;type:varchar(255)" json:"role"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *Role) TableName() string {
	return "role"
}

type Admin struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key" json:"id"`
	RoleId      int64     `gorm:"column:role_id;type:bigint(20)" json:"role_id"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Password    string    `gorm:"column:password;type:varchar(255)" json:"password"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
}

func (m *Admin) TableName() string {
	return "admin"
}

type Orders struct {
	Id          int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	OrderNo     string    `gorm:"column:order_no;type:varchar(255)" json:"order_no"`
	Status      string    `gorm:"column:status;type:varchar(255)" json:"status"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
	CreatedTime time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
}

func (m *Orders) TableName() string {
	return "orders"
}
