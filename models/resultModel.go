package models

import (
	"time"
)

type Users struct {
	ID         uint
	Created_at time.Time
	Updated_at time.Time
	UserName   string
	Email      string
	FName      string
	LName      string
	City       string
	Birth      string
	Gender     string
	Phone      string
	Address    string
	Role       string
	Score      int
}
type Categories struct {
	ID   uint
	Name string
	Description string
}

type Books struct {
	Id          uint
	Title       string
	Price       int	
	Author      string
	Isbc        string
	Edition     string
	Pages       string
	Year        string
	Language    string
	Sabject     string
	Publisher   string
	Translator  string
}

type BookInfo struct {
	Id          uint
	Rating      float64
	Title       string
	Author      string
	Isbc        string
	Edition     string
	Pages       string
	Year        string
	Language    string
	Sabject     string
	Publisher   string
	Translator  string
	About       string
}

type Quantity_books struct {
	Id    uint
	Updated_at time.Time
	Created_at time.Time
	Quantity int
	SPrice int
}

type Quantities struct {
	Id    uint
	Quantity int
	Book_id uint
	BPrice int
	SPrice int
	Updated_at time.Time
	Created_at time.Time
}

type Discounts struct {
	Id      uint
	Book_id uint
	Discount int
	Deadline string
}

type Comments struct {
	Id         uint
	Title      string
	Book_id    uint
	Comment    string
	Replayid   uint
	Status     string
	Created_at time.Time
	Likes      int
	Dislikes   int
}

type Comment_bookss struct {
	Id         uint
	User_id    uint
	User_name  string
	created_at time.Time
	Replayid    uint
	Comment     string
	Likes      int
	Dislikes   int
}

type CheckComment struct {
	Id         uint
	Created_at time.Time
	Comment    string
	Replayid   uint
	User_id    uint
	User_name  string
	Book_id    uint
	Title      string
}

type Orders struct {
	Id         uint
	User_id    uint
	Status     string
	Created_at time.Time
}
type OrderInfo struct {
	QBook_id uint
	Quantity int
	SPrice   int
	Book_id  uint
	Title    string
}