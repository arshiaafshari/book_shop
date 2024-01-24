package controllers

import (
	"example/user/initializers"
	"example/user/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//books
func AddBook(c *gin.Context) {
	// get the information
	var body struct {
		Title      string
		Author     string
		ISBC       string
		Edition    string
		Pages      string
		Year       string
		Language   string
		Sabject    string
		Publisher  string
		Translator string
		CategoryId []uint
		About      string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}

	// create the book 

	book := models.Book{
		Title:      body.Title,
		Author:     body.Author,
		ISBC:       body.ISBC,
		Edition:    body.Edition,
		Pages:      body.Pages,
		Year:       body.Year,
		Language:   body.Language,
		Sabject:    body.Sabject,
		Publisher:  body.Publisher,
		Translator: body.Translator,
	}
	tx := initializers.DB.Begin()
	//insert the book
	result := tx.Create(&book)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to add book",
		})
		tx.Rollback()
		return
	}
	//create the about
	about := models.About_book{
		About: body.About,
		Book_id: book.ID,
	}
	//insert the about
	rgorm := tx.Create(&about)

	if rgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to add about book",
		})
		tx.Rollback()
		return
	}
	//create the category
	for i := 0; i < len(body.CategoryId); i++ {
		
	category := models.Category_book{
		Category_id: body.CategoryId[i],
		Book_id:     book.ID,
	}
	// insert the category
	cgorm := tx.Create(&category)
	if cgorm.Error != nil {
	    c.JSON(http.StatusBadRequest, gin.H{
	        "error": "Failed to add category",
	    })
	    tx.Rollback()
	    return
	}
    }
	//respond
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Add book success",
	})
}

func UpdateBook(c *gin.Context) {
	var body struct {
		ID         uint
		Title      string
		Author     string
		ISBC       string
		Edition    string
		Pages      string
		Year       string
		Language   string
		Sabject    string
		Publisher  string
		Translator string
		CategoryId []uint
		About      string

	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find the book and about book
	var book models.Book
	initializers.DB.First(&book, "id = ?", body.ID)
	
	if book.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't find book",
		})
		return
	}

	var about models.About_book
	initializers.DB.First(&about, "book_id = ?", body.ID)

	if about.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Can't find about book",
		})
		return
	}
	//find the category book
	var category []models.Category_book
	var category_id []uint
	initializers.DB.Where("book_id = ?", body.ID).Find(&category)

	for i := 0; i < len(category); i++ {
		category_id = append(category_id, category[i].Category_id)
	}
	tx := initializers.DB.Begin()
	//update
	if body.Title != book.Title ||
	body.Author != book.Author ||
	body.ISBC != book.ISBC ||
	body.Edition != book.Edition ||
	body.Pages != book.Pages ||
	body.Year != book.Year ||
	body.Language != book.Language ||
	body.Sabject != book.Sabject ||
	body.Publisher != book.Publisher ||
	body.Translator != book.Translator {	
		dgorm := tx.Model(&book).Updates(models.Book{
		Title:      body.Title,
		Author:     body.Author,
		ISBC:       body.ISBC,
		Edition:    body.Edition,
		Pages:      body.Pages,
		Year:       body.Year,
		Language:   body.Language,
		Sabject:    body.Sabject,
		Publisher:  body.Publisher,
		Translator: body.Translator,
		})
		if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update book",
		})
		tx.Rollback()
		return
		}
    }
	//update about
	if body.About != about.About {
		ugorm := initializers.DB.Model(&about).Updates(models.About_book{
			About: body.About,
		})
		if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update about book",
		})
		tx.Rollback()
		return
		}
    }
	//update category

	for i := 0; i < len(body.CategoryId); i++ {
		for j := 0; j < len(category_id); j++ {
			if body.CategoryId[i] == category_id[j] {
				//remove the category from the category_id
				category_id = append(category_id[:j], category_id[j+1:]...)
				break
			}
			if j == len(category_id) - 1 {
				category := models.Category_book{
					Category_id: body.CategoryId[i],
					Book_id:     book.ID,
				}
				// insert the category
				cgorm := tx.Create(&category)
				if cgorm.Error != nil {
				    c.JSON(http.StatusBadRequest, gin.H{
				        "error": "Failed to add category",
				    })
				    tx.Rollback()
				    return
				}
			}
		}
	}
	//delete the category
	for i := 0; i < len(category_id); i++ {
		dgorm := tx.Where("book_id = ? and category_id = ?", book.ID, category_id[i]).Delete(&models.Category_book{})
		if dgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to delete category",
			})
			tx.Rollback()
			return
		}
	}
	//respond
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "book update success",
	})
}

func DeleteBook(c *gin.Context) {
	var body struct {
		Id uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	BookId := body.Id
	//get the book
	var book models.Book

	initializers.DB.First(&book, "id = ?", BookId)

	if book.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})

		return
	}

	//soft delete
	dgorm :=  initializers.DB.Delete(&book)

	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete book",
		})
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "book delete success",
	})
}

func GetBook(c *gin.Context) { 

	var books []models.Books	
	qgorm := initializers.DB.Raw(  `SELECT "id", title, t1.price , author, publisher, "year", "language", isbc, edition, pages, translator 
									FROM books
									LEFT JOIN (	SELECT quantity_books.book_id , Min(quantity_books.s_price) AS price
												FROM quantity_books
												WHERE quantity_books.quantity > 0 AND quantity_books.deleted_at IS NULL
												GROUP BY quantity_books.book_id) AS t1 ON books."id" = t1.book_id
									WHERE books.deleted_at IS NULL`).Scan(&books)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get book",
		})
		return
	}
	
	c.JSON(http.StatusOK, books)
}

func GetBookInfo(c *gin.Context) {

	//get id
	BookId := c.Param("id")
	
	//find the book
	var info models.BookInfo
	qgorm := initializers.DB.Raw(`	SELECT *
									FROM (	SELECT rating_books.book_id  , AVG(rating_books.rating) AS rating
											FROM rating_books
											WHERE rating_books.book_id = ? AND rating_books.deleted_at IS NULL
											GROUP BY rating_books.book_id) AS t1
									RIGHT JOIN (SELECT books."id", books.title , books.author, books.publisher , books."year" , books."language" , books.isbc , books.edition , books.pages            ,books.translator , about_books.about
												FROM books
												JOIN about_books ON books."id" = about_books.book_id
												WHERE books."id" = ? AND books.deleted_at IS NULL AND about_books.deleted_at IS NULL) AS t2 ON t1.book_id = t2."id"`, BookId, BookId).Scan(&info) 

	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get book",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, info)
}

func GetBookPrice(c *gin.Context) {
	//get id
	BookId := c.Param("id")
	
	//find quntity and prices
	var qprice []models.Quantity_books

	qgorm := initializers.DB.Find(&qprice, "book_id = ? AND deleted_at IS NULL", BookId)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get qprice",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, qprice)
}

func GetBookComment(c *gin.Context) {
	//get id
	BookId := c.Param("id")
	
	//find comment
	var comment []models.Comment_bookss

	qgorm := initializers.DB.Raw(`	SELECT t1.user_id , t1.user_name , t1."id" , t1.created_at , t1.replayid , t1."comment" , t2.likes , t2.dislikes
									FROM (	SELECT comment_books.user_id , users.user_name , comment_books."id" , comment_books.created_at , comment_books.replayid , comment_books."comment"
											FROM comment_books
											LEFT JOIN users ON comment_books.user_id = users."id"
											WHERE comment_books.book_id = ? AND comment_books.deleted_at IS NULL AND users.deleted_at IS NULL) AS t1
									LEFT JOIN ( SELECT cl1.comment_id 
													,  COUNT(CASE WHEN cl1.like IS TRUE THEN cl1."like" END) AS likes 
													,  COUNT(CASE WHEN cl1.like IS FALSE THEN  cl1."like" END) AS dislikes
												FROM comment_likes AS cl1
												WHERE cl1.deleted_at IS NULL
												GROUP BY cl1.comment_id) AS t2 ON t1."id" = t2.comment_id`,  BookId).Scan(&comment)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get comments",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, comment)
}

func GetBookCategory(c *gin.Context) {

	//get id
	BookId := c.Param("id")
	fmt.Println(BookId)
	//find category
	var category []models.Categories
	qgorm := initializers.DB.Joins("JOIN category_books ON category_books.category_id = categories.id").Where("category_books.book_id = ? and category_books.deleted_at IS NULL", BookId).Find(&category)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get category",
		})
		
		return
	}
	//respond
	c.JSON(http.StatusOK, category)
}

func AddRating(c *gin.Context) {

	var body struct {
		ID uint
		Rating int
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	BookId := body.ID
	Rating := body.Rating
	//check not repit
	var rating models.Rating_book
	initializers.DB.First(&rating, "book_id = ? AND user_id = ?", BookId, c.GetUint("USERID"))
	tx := initializers.DB.Begin()
	if rating.ID != 0 {	
		// delete
		dgorm := tx.Delete(&rating)
		if dgorm.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to delete last rating",
			})
			return	
		}
	}
	//create
	qgorm := tx.Create(&models.Rating_book{
		Book_id: BookId,
		Rating: Rating,
		User_id: c.GetUint("USERID"),
	})
	if qgorm.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create rating",
		})
		return
	}
	//respond
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "rating created",
	})
}

//category
func GetCategory(c *gin.Context) {
	var category []models.Category
	
	qgorm := initializers.DB.Select("id, name").Find(&category)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get category",
		})
		return
	}
	var categories []models.Categories
	for i := 0; i < len(category); i++ {
		categories = append(categories, models.Categories{
			ID:    category[i].ID,
			Name:  category[i].Name,
		})
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategory_book(c *gin.Context) {
	var category_book []models.Category_book
	
	qgorm := initializers.DB.Find(&category_book)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get category_book",
		})
		return
	}

	c.JSON(http.StatusOK, category_book)
}

func AddCategory(c *gin.Context) {
	var body struct {
		Name string
		Discription string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//check name
	var category models.Category
	initializers.DB.First(&category, "name = ?", body.Name)
	if category.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Category already exist",
		})
		return
	}
	//create
	qgorm := initializers.DB.Create(&models.Category{
		Name:       body.Name,
		Description: body.Discription,
	})
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create category",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "category create success",
	})
}

func UpdateCategory(c *gin.Context) {
	var body struct {
		ID uint
		Name string
		Discription string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find category
	var category models.Category
	qgorm := initializers.DB.First(&category, "id = ?", body.ID)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Category not found",
		})
		return
	}
	//update
	dgorm := initializers.DB.Model(&category).Updates(models.Category{
		Name:       body.Name,
		Description: body.Discription,
	})
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update category",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "category update success",
	})
}

func DeleteCategory(c *gin.Context) {
	var body struct {
		ID uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	fmt.Println(body.ID)
	//find category
	var category models.Category
	initializers.DB.First(&category, "id = ?", body.ID)
	if category.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Category not found",
		})
		return
	}
	//soft delete
	dgorm :=  initializers.DB.Delete(&category)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete category",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "category delete success",
	})		
}

//quantity
func AddQuantity(c *gin.Context) {
	var body struct {
		Idbook  uint
		Count    int
		Bprice   int
		Sprice   int
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find book title
	var book models.Book
	qqgorm := initializers.DB.First(&book, "id = ?", body.Idbook)
	if qqgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find book",
		})
		return
	}	
	//create
	tx := initializers.DB.Begin()
	qgorm := tx.Create(&models.Quantity_book{
		Book_id:  body.Idbook,
		Quantity: body.Count,
		BPrice:   body.Bprice,
		SPrice:   body.Sprice,
	})
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create quantity",
		})
		return
	}
	//check bookmark
	if body.Count > 0 {
	var bookmark []models.UserBookmark
	dgorm := initializers.DB.Select("user_id").Find(&bookmark, "book_id = ?", body.Idbook)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to check bookmark",
		})
		tx.Rollback()
		return
	}

	//build notification
	for i := 0; i < len(bookmark); i++ {
		qgorm := initializers.DB.Create(&models.UserNotification{
			User_id:      bookmark[i].User_id,
			Notification: book.Title + "is available now ",
			Seen:        false,
		})
		if qgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create notification",
			})
			tx.Rollback()
			return
		}
	}

	//delete bookmark
	for i := 0; i < len(bookmark); i++ {
		dgorm := initializers.DB.Delete(&models.UserBookmark{}, "user_id = ? AND book_id = ?", bookmark[i].User_id, body.Idbook)
		if dgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to delete bookmark",
			})
			tx.Rollback()
			return
		}
	}
    }
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "quantity create success",
	})
}

func UpdateQuantity(c *gin.Context) {
	var body struct {
		ID       uint
		Idbook 	 uint
		Count 	 int
		Bprice   int
		Sprice   int
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find quantity
	var quantity models.Quantity_book
	qgorm := initializers.DB.First(&quantity, "id = ?", body.ID)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quantity not found",
		})
		return
	}
	//update
	dgorm := initializers.DB.Model(&quantity).Updates(models.Quantity_book{
		Book_id:  body.Idbook,
		Quantity: body.Count,
		BPrice:   body.Bprice,
		SPrice:   body.Sprice,
	})
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update quantity",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "quantity update success",
	})
}

func DeleteQuantity(c *gin.Context) {
	//get id
	var body struct {
		ID uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	quantityId := body.ID
	//find quantity
	var quantity models.Quantity_book
	qgorm := initializers.DB.First(&quantity, "id = ?", quantityId)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quantity not found",
		})
		return
	}
	//soft delete
	dgorm :=  initializers.DB.Delete(&quantity)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete quantity",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "quantity delete success",
	})
}

func GetQuantity(c *gin.Context) {
	var quantity []models.Quantity_book
	qgorm := initializers.DB.Find(&quantity)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quantity not found",
		})
		return
	}
	var quantities []models.Quantities
	for i := 0; i < len(quantity); i++ {
		quantities = append(quantities, models.Quantities{
			Id: quantity[i].ID,
			Quantity: quantity[i].Quantity,
			Book_id:  quantity[i].Book_id,
			BPrice:   quantity[i].BPrice,
			SPrice:   quantity[i].SPrice,
			Updated_at: quantity[i].UpdatedAt,
		})
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"quantity": quantities,
	})
}

//discount
func AddDiscount(c *gin.Context) {
	UserId := c.GetUint("USERID")
	var  body struct { 
		Bookid  uint
		Discount int
		Startdate string
		Deadline string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//create
	qgorm := initializers.DB.Create(&models.Discount_book{
		Admin_id : UserId,
		Book_id:  body.Bookid,
		Discount: body.Discount,
		StartDate: body.Startdate,
		Deadline: body.Deadline,
	})
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create discount",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "discount create success",
	})
}

func UpdateDiscount(c *gin.Context) {
	UserId := c.GetUint("USERID")

	var  body struct { 
		ID       uint
		Bookid  uint
		Discount int
		Startdate string
		Deadline string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find discount
	var discount models.Discount_book
	qgorm := initializers.DB.First(&discount, "id = ?", body.ID)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Discount not found",
		})
		return
	}
	//update
	dgorm := initializers.DB.Model(&discount).Updates(models.Discount_book{
		Admin_id : UserId,
		Book_id:  body.Bookid,
		Discount: body.Discount,
		StartDate: body.Startdate,
		Deadline: body.Deadline,
	})
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update discount",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "discount update success",
	})
}

func DeleteDiscount(c *gin.Context) {
	//get id
	var body struct {
		ID uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	discountId := body.ID
	//find discount
	var discount models.Discount_book
	qgorm := initializers.DB.First(&discount, "id = ?", discountId)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Discount not found",
		})
		return
	}
	//soft delete
	dgorm :=  initializers.DB.Delete(&discount)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete discount",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "discount delete success",
	})
}

func GetDiscount(c *gin.Context) {
	currentTime := time.Now()
	date :=  strconv.Itoa(int(currentTime.Year()))+"-"+strconv.Itoa(int(currentTime.Month()))+"-"+strconv.Itoa(int(currentTime.Day()))
	fmt.Println(date)
	var discounts []models.Discount_book
	qgorm := initializers.DB.Where("deadline >= ? and start_date <= ?", date, date).Find(&discounts)
	var discounts2 []models.Discounts
	for i := 0; i < len(discounts); i++ {
		discounts2 = append(discounts2, models.Discounts{
			Id: discounts[i].ID,
			Book_id: discounts[i].Book_id,
			Discount: discounts[i].Discount,
			Deadline: discounts[i].Deadline,
		})
	}
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Discount not found",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, discounts2)
}


//comment
func AddComment(c *gin.Context) {
	var body struct {
		Comment  string
		Bookid   uint
		Replayid uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	UserId := c.GetUint("USERID")
	BookId := body.Bookid
	comment := body.Comment
	replayid := body.Replayid
	//find replay
	if replayid != 0 {
		var comment2 models.Comment_book
		initializers.DB.First(&comment2, "id = ? and status = ?", replayid, "confirmation")
		if comment2.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Replaid comment not found",
			})
			return
		}
	}
	//create
	qgorm := initializers.DB.Create(&models.Comment_book{
		User_id: UserId,
		Book_id: BookId,
		Comment: comment,
		Replayid: replayid,
		Status: "not confirmed",
	})
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create comment",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "comment create success",
	})
}

func DeleteComment(c *gin.Context) {
	var body struct {
		ID uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	CommentId := body.ID
	//find comment
	var comment models.Comment_book
	initializers.DB.First(&comment, "id = ?", CommentId)
	if comment.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Comment not found",
		})
		return
	}
	//check permision
	if (comment.User_id != c.GetUint("USERID") && c.GetString("Role") == "user") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Don't have permision to delete this comment",
		})
		return
	}
	//soft delete
	qgorm := initializers.DB.Delete(&comment)

	if qgorm.Error != nil {
	    c.JSON(http.StatusBadRequest, gin.H{
	        "error": "Failed to delete comment",
	    })
	    return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "comment delete success",
	})
}

func MyComments(c *gin.Context) {
	var comments []models.Comments
	qgorm := initializers.DB.Raw(`SELECT ids AS id , title , book_id , comment , replayid , status , created_at , l1.likes , l1.dislikes
								  FROM (SELECT comment_books.id AS ids , title , book_id , comment , replayid , status , comment_books.created_at AS created_at
										FROM books
										JOIN comment_books ON comment_books.book_id = books."id"
										WHERE comment_books.deleted_at IS NULL AND comment_books.user_id = ?)
							      LEFT JOIN	( SELECT  cl1.comment_id 
													, COUNT(CASE WHEN cl1.like IS TRUE THEN cl1."like" END) AS likes 
													, COUNT(CASE WHEN cl1.like IS FALSE THEN  cl1."like" END) AS dislikes
											  FROM comment_likes AS cl1
											  WHERE cl1.deleted_at IS NULL
											  GROUP BY cl1.comment_id) AS l1  ON ids = l1.comment_id
								  ORDER BY created_at DESC`	, c.GetUint("USERID")).Scan(&comments)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Comment not found",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, comments)
}

func CheckComment(c *gin.Context) {
	var comments []models.CheckComment
	qgorm := initializers.DB.Raw("SELECT cb.id AS id , cb.created_at AS created_at , cb.comment AS comment , cb.replayid AS replayid , cb.user_id AS user_id , users.user_name AS user_name , cb.book_id AS book_id , cb.title AS title FROM ( SELECT comment_books.id , comment_books.created_at , comment_books.comment , comment_books.replayid , comment_books.user_id , comment_books.book_id , books.title FROM comment_books LEFT JOIN books ON comment_books.book_id = books.id WHERE comment_books.status = 'not confirmed' ) AS cb LEFT JOIN users ON cb.user_id = users.id").Scan(&comments)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Comments can't loaded",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, comments)
}

func ConfirmComment(c *gin.Context) {
	var body struct {
		ID      uint
		Confirm bool
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	CommentId := body.ID
	var status string
	if body.Confirm {
		status = "confirmation"
	}else{
		status = "blocked"
	}
	//find comment
	var comment models.Comment_book
	initializers.DB.First(&comment, "id = ?", CommentId)
	if comment.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Comment not found",
		})
		return
	}
	//update
	dgorm := initializers.DB.Model(&comment).Updates(models.Comment_book{
		Status: status,
	})
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to change status",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "comment update status success",
	})
}

func LikeComment(c *gin.Context) {
	var body struct {
		Like bool
		Commentid uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//check comment exist
	var comment models.Comment_book
	initializers.DB.First(&comment, "id = ? AND status = ?", body.Commentid, "confirmation")
	if comment.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Comment not found",
		})
		return
	}
	//find a like with this user
	var like models.Comment_like
	initializers.DB.First(&like, "user_id = ? AND comment_id = ?", c.GetUint("USERID"), body.Commentid)
	//if any like for this user on this comment 
	if like.ID == 0 {
		//create
		igorm := initializers.DB.Create(&models.Comment_like{
			User_id: c.GetUint("USERID"),
			Comment_id: body.Commentid,
			Like: body.Like,
		})
		if igorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to like comment",
			})
			return
		}
	}else{
	//update
	if like.Like != body.Like {
	ugorm := initializers.DB.Model(&like).Update("like", body.Like)
	if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to like comment",
		})
		return
	}
	}else{
		if body.Like {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "You already like this comment",
			})
			return
		}else{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "You already dislike this comment",
			})
			return
		}
	}
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "like or dislike success",
	})
}

//bookmark
func AddBookmark(c *gin.Context) {
	var body struct{
		Id uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to read body",
		})
		return
	}
	var bookmark models.UserBookmark 
	initializers.DB.First(&bookmark, "user_id = ? AND book_id = ?", c.GetUint("USERID"), body.Id)
	if bookmark.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "This bookmark already exist!",
		})
		return
	}
	//check book quntity exist
	var quntity models.Quantity_book
	initializers.DB.First(&quntity, "book_id = ? AND quantity > 0", body.Id)
	if quntity.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "This book is available don't need add bookmark!",
		})
	}
	//create
	igorm := initializers.DB.Create(&models.UserBookmark{
		User_id: c.GetUint("USERID"),
		Book_id: body.Id,
	})
	if igorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add bookmark",
		})
		return	
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "add bookmark success",
	})
}

func DeleteBookmark(c *gin.Context) {
	var body struct{
		Id uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to read body",
		})
		return
	}
    //find bookmark
	var bookmark models.UserBookmark 
	initializers.DB.First(&bookmark, "id = ?", body.Id)
	if bookmark.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "This bookmark not exist!",
		})
		return
	}
	//check 
	if bookmark.User_id != c.GetUint("USERID") && c.GetString("Role") == "user" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "You can't delete this bookmark!",
		})
		return
		
	}
	//delete
	dgorm := initializers.DB.Delete(&bookmark)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete bookmark",
		})
		return	
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "delete bookmark success",
	})
}

func GetBookmarks(c *gin.Context) {
	var bookmarks []models.UserBookmark
	//find user's bookmarks
	if c.GetString("Role") == "user" {
	qgorm :=  initializers.DB.Where("user_id = ?", c.GetUint("USERID")).Find(&bookmarks)
		if qgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bookmarks not found",
			})
			return
		}
	}else{
		//find all bookmarks
		qgorm :=  initializers.DB.Find(&bookmarks)
		if qgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bookmarks not found",
			})
			return
		}
	}
	//respond
	c.JSON(http.StatusOK, bookmarks)
}


//notifications
func NewNotifications(c *gin.Context) {
	var notifications []models.UserNotification
	//find user's notifications
	qgorm :=  initializers.DB.Where("user_id = ? AND seen = false", c.GetUint("USERID")).Find(&notifications)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Notifications not found",
		})
		return
	}
	//update seens
	ugorm :=  initializers.DB.Model(&notifications).Update("seen", true)
	if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update notifications seen",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, notifications)
}

func GetNotifications(c *gin.Context) {
	var notifications []models.UserNotification
	//find user's notifications
	qgorm :=  initializers.DB.Where("user_id = ?", c.GetUint("USERID")).Find(&notifications)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Notifications not found",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, notifications)
}