package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"unique"`
	Password string `gorm:"not null"`
	Email    string 
	FName    string
	LName    string 
	City     string 
	Birth    string `gorm:"type:date ; default:0001-01-01"`
	Gender   string `gorm:"type: gender_enum ; default:'male' ;not null"`
	Phone    string `gorm:"like:'[0-9]{10}'"`
	Address  string
	Role     string `gorm:"type: role_enum; default:'user' ;not null"`
	Jwt_id   string 
	Score    int    `gorm:"default:0"`

	Comment_books []Comment_book `gorm:"foreignKey:User_id"`
	Rating_books []Rating_book `gorm:"foreignKey:User_id"`
	Discount_books []Discount_book `gorm:"foreignKey:Admin_id"`
	Orders []Order `gorm:"foreignKey:User_id"`
	UserBookmarks []UserBookmark `gorm:"foreignKey:User_id"`
	UserNotifications []UserNotification `gorm:"foreignKey:User_id"`
	Comment_likes []Comment_like `gorm:"foreignKey:User_id"`
}

type Book struct {
	gorm.Model
	Title      string `gorm:"unique ;not null"`
	Author     string 
	ISBC       string `gorm:"unique ;not null"`
	Edition    string 
	Pages      string
	Year       string 
	Language   string 
	Img    	   bool   `gorm:"default:false"`
	Sabject    string 
	Publisher  string 
	Translator string 

	About_books []About_book `gorm:"foreignKey:Book_id"`
	Comment_books []Comment_book `gorm:"foreignKey:Book_id"`
	Rating_books []Rating_book `gorm:"foreignKey:Book_id"`
	Quantity_books []Quantity_book `gorm:"foreignKey:Book_id"`
	Category_books []Category_book `gorm:"foreignKey:Book_id"`
	UserBookmarks []UserBookmark `gorm:"foreignKey:Book_id"`
}

type About_book struct {
	gorm.Model
	Book_id     uint 
	About       string 
}

type Comment_book struct {
	gorm.Model
	Book_id     uint 
	Comment     string
	Replayid    uint      
	User_id     uint 
	Status      string `gorm:"type: comment_status_enum ; default:'not confirmed' ;not null"`
	Comments []Comment_like `gorm:"foreignKey:Comment_id"`
}

type Comment_like struct {
	gorm.Model
	Comment_id  uint
	User_id     uint
	Like        bool `gorm:"not null"`
}

type Rating_book struct {
	gorm.Model
	Book_id     uint 
	Rating      int	`gorm:"between:1,5"`
	User_id     uint
}

type Discount_book struct {
	gorm.Model
	Admin_id    uint
	Book_id     uint `gorm:"foreignKey:"`
	Discount     int `gorm:"between:1,100"`
	StartDate   string `gorm:"type:date"`
	Deadline    string `gorm:"type:date"`
}

type Category struct {
	gorm.Model
	Name        string `gorm:"unique ;not null"`
	Description string

	Category_books []Category_book `gorm:"foreignKey:Category_id"`
}

type Category_book struct {
	gorm.Model
	Book_id     uint 
	Category_id uint 
}

type Order struct {
	gorm.Model
	User_id     uint 
	Status      string `gorm:"type: order_status_enum ; default:'not Finished' ;not null"`
//	LStatus      string `gorm:"type: order_status_enum ; default:'not Finished' ;not null"`
	Order_books []Order_book `gorm:"foreignKey:Order_id"`
}

type Quantity_book struct {
	gorm.Model
	Book_id     uint 
	BPrice      int
	SPrice      int
	Quantity    int `gorm:"unsigned"`
	Order_books []Order_book `gorm:"foreignKey:QBook_id"`
}

type Order_book struct {
	gorm.Model
	Order_id uint 
	QBook_id  uint 
	Quantity int `gorm:"unsigned"`
}

type  UserBookmark struct {
	gorm.Model
	Book_id     uint 
	User_id     uint 
}

type UserNotification struct {
	gorm.Model
	User_id     uint
	Notification string
	Seen        bool `gorm:"default:false"`
}

