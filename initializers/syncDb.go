package initializers

import "example/user/models"

func SyncDb() {
	DB.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.Category{},
		&models.Category_book{},
		&models.Order{},
		&models.Order_book{},
		&models.About_book{},
		&models.Rating_book{},
		&models.Comment_book{},
		&models.Discount_book{},
		&models.UserBookmark{},
		&models.UserNotification{},
		&models.Comment_like{},
		&models.Quantity_book{})
}
