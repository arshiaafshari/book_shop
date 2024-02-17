package controllers

import (
	"example/user/initializers"
	"example/user/models"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetOrder(c *gin.Context) {
	var body struct {
		UserId uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if c.GetString("role") == "user" && c.GetUint("USERID") != body.UserId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not your order",
		})
		return
	}
	qgorm := initializers.DB.Create(&models.Order{
		User_id: body.UserId,
		Status:  "not Finished"})
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to create order",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "order created",
	})
}

func AddToOrder(c *gin.Context) {
	var body struct {
		Order_id uint
		QBook_id uint
		Quantity int
		Discount_id uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//check order
	var order models.Order
	initializers.DB.First(&order, "id = ?", body.Order_id)
	if order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	if order.User_id != c.GetUint("USERID") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not your order",
		})
		return
	}
	if order.Status != "not Finished" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot add to this order",
		})
		return
	}
	//check book and quantity
	var quntity models.Quantity_book
	initializers.DB.First(&quntity, "id = ?", body.QBook_id)
	if quntity.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "book not found",
		})
		return
	}
	//create or update
	var order_book models.Order_book
	initializers.DB.Where("order_id = ? and q_book_id = ?", body.Order_id, body.QBook_id).First(&order_book)
	if order_book.ID == 0 {
		qgorm := initializers.DB.Create(&models.Order_book{
			Order_id: body.Order_id,
			QBook_id: body.QBook_id,
			Quantity: body.Quantity,
		})
		if qgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create quantity",
			})
			return
		}
	} else {
		qgorm := initializers.DB.Model(&order_book).Update("quantity", body.Quantity)
		if qgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to update quantity",
			})
			return
		}
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "add to order success",
	})
}

func DeleteFromOrder(c *gin.Context) {
	//get body
	var body struct {
		Orb_id uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find
	var order_book models.Order_book
	initializers.DB.First(&order_book, "id = ?", body.Orb_id)
	if order_book.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "this pat of order not found",
		})
		return
	}
	var order models.Order
	initializers.DB.First(&order, "id = ?", order_book.Order_id)
	if order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	if order.User_id != c.GetUint("USERID") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not your order",
		})
		return
	}
	//delete
	dgorm := initializers.DB.Delete(&order_book)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete quantity",
		})
		return
	}
	//resond
	c.JSON(http.StatusOK, gin.H{
		"message": "this book deleted from order",
	})
}

func DeleteOrder(c *gin.Context) {
	//Get order id
	var body struct {
		OrderId uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	OrderId := body.OrderId
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ?", OrderId)
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	if order.User_id != c.GetUint("USERID") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not yours",
		})
		return
	}
	if order.Status != "not Finished" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot delete this order",
		})
		return
	}
	//delete
	dgorm := initializers.DB.Delete(&order)
	if dgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete order",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "order delete success",
	})
}

func FinishOrder(c *gin.Context) {
	//Get order_id
	var body struct {
		OrderId uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ? ", body.OrderId)
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	} else {
		if c.GetString("role") == "user" && order.User_id != c.GetUint("USERID") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Not your order",
			})
			return
		}
		tx := initializers.DB.Begin()
		//check and update quantity
		var orderbook []models.Order_book
		var quntity models.Quantity_book
		tx.Where("order_id = ? ", order.ID).Find(&orderbook)
		for i := 0; i < len(orderbook); i++ {
			tx.Where("id = ? ", orderbook[i].QBook_id).First(&quntity)
			if quntity.Quantity < orderbook[i].Quantity {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Not enough " + strconv.Itoa(int(quntity.ID)),
				})
				return
			} else {
				orm := tx.Model(&quntity).Update("quantity", quntity.Quantity-orderbook[i].Quantity)
				if orm.Error != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Failed to update quantity",
					})
					return
				}
			}
		}

		//update
		ugorm := tx.Model(&order).Update("status", "waiting for confirmation")
		if ugorm.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed go to Finshed status",
			})
			return
		}
		tx.Commit()
	}

	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "go to Finshed status success",
	})
}

func ConfirmOrder(c *gin.Context) {
	//Get order id and confirm bool
	var body struct {
		OrderId uint
		Confirm bool
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	OrderId := body.OrderId
	Confirm := body.Confirm
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ? AND status = 'waiting for confirmation'", OrderId)
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	//quantification status
	var Status string
	if Confirm {
		Status = "confirmed"
	} else {
		Status = "rejected"
	}
	//update
	ugorm := initializers.DB.Model(&order).Update("status", Status)
	if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to confirm order",
		})
		return
	}
	//set a notification
	cgorm := initializers.DB.Create(&models.UserNotification{
		User_id:      order.User_id,
		Notification: "order has been " + Status,
	})
	if cgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create notification",
		})
	}
	//respond

	c.JSON(http.StatusOK, gin.H{
		"message": "order confirm success",
	})
}

func UpdateOrder(c *gin.Context) {
	//Get order id
	OrderId := c.GetUint("order_id")

	UserId := c.GetUint("user_id")
	Status := c.GetString("status")
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ?", OrderId)
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	//update
	ugorm := initializers.DB.Model(&order).Update("user_id", UserId).Update("l_status", order.Status).Update("status", Status)
	if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update order",
		})
		return
	}
	//respond

	c.JSON(http.StatusOK, gin.H{
		"message": "order update success",
	})
}

func CancelOrder(c *gin.Context) {
	//Get order id
	OrderId := c.GetUint("order_id")
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ? AND user_id = ? ", OrderId, c.GetUint("USERID"))
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	if order.Status == "canceled" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order already canceled",
		})
		return
	}
	if order.Status == "wating for canselation" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your canselation request is already sent, please wait for the confirmation",
		})
		return
	}
	if order.Status == "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your order rejected , cannot be cancel",
		})
		return
	}
	//update
	ugorm := initializers.DB.Model(&order).Update("status", "wating for canselation")
	if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to cancele order",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "cancel request send",
	})
}

func Canselation(c *gin.Context) {
	//Get order id
	OrderId := c.GetUint("order_id")
	Cansel := c.GetBool("cansel")
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ? AND status = 'wating for canselation' ", OrderId)
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	//qunatification status
	var Status string
	var note string
	if Cansel {
		Status = "canceled"
		note = "Your order has been canceled"
	} else {
		Status = order.Status
		note = "Your order has been not canceled and back to last status"
	}
	//update
	ugorm := initializers.DB.Model(&order).Update("l_status", order.Status).Update("status", Status)
	if ugorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to cancele order",
		})
		return
	}
	//set a notification
	cgorm := initializers.DB.Create(&models.UserNotification{
		User_id:      order.User_id,
		Notification: note,
	})
	if cgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create notification",
		})
	}
	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "cancelation success",
	})
}

func GetOrders(c *gin.Context) {

	var orders []models.Orders
	//check user
	if c.GetString("Role") == "user" {
		initializers.DB.Find(&orders, "user_id = ? ", c.GetUint("USERID"))
	} else {
		initializers.DB.Find(&orders)
	}
	//respond
	c.JSON(http.StatusOK, orders)
}

func GetWaitingOrders(c *gin.Context) {
	var orders []models.Orders
	initializers.DB.Find(&orders, "status = 'waiting for cancelation' OR status = 'waiting for confirmation'")
	//respond
	c.JSON(http.StatusOK, orders)
}

func GetOrderInfo(c *gin.Context) {
	//Get order id
	OrderId := c.Param("id")
	//find and check order
	var order models.Order
	ogorm := initializers.DB.First(&order, "id = ? ", OrderId)
	if ogorm.Error != nil || order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	if c.GetString("Role") == "user" && order.User_id != c.GetUint("USERID") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Don't have permision to see this order",
		})
		return
	}
	//find info
	var orderinfo []models.OrderInfo
	qgorm := initializers.DB.Raw(`SELECT  t1.q_book_id , t1.quantity , t1.book_id , t1.s_price , books.title
								  FROM	  (SELECT order_books.q_book_id , order_books.quantity , quantity_books.book_id , quantity_books.s_price
				 						   FROM order_books
										   LEFT JOIN quantity_books ON order_books.q_book_id = quantity_books."id"
				 						   WHERE order_books.order_id= ? AND order_books.deleted_at IS NULL ) AS t1
								  LEFT JOIN books ON 	books."id" = t1.book_id`, OrderId).Scan(&orderinfo)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get order info",
		})
		return
	}
	//respond
	c.JSON(http.StatusOK, orderinfo)
}

func UserChangeOrderStatus(c *gin.Context) {
	var body struct {
		Status string
		Id     uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find order
	var order models.Order
	initializers.DB.First(&order, "id = ? ", body.Id)
	if order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}
	//check user
	if order.User_id != c.GetUint("USERID") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Don't have permision to change order status",
		})
		return
	}
	//update
	if body.Status == "waiting for confirmation" {
		if order.Status != "not Finished" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't change order status",
			})
			return
		}

		tx := initializers.DB.Begin()
		//check and update quantity
		var orderbook []models.Order_book
		var quntity models.Quantity_book
		tx.Where("order_id = ? ", order.ID).Find(&orderbook)
		for i := 0; i < len(orderbook); i++ {
			tx.Where("id = ? ", orderbook[i].QBook_id).First(&quntity)
			if quntity.Quantity < orderbook[i].Quantity {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Not enough " + strconv.Itoa(int(quntity.ID)),
				})
				return
			} else {
				orm := tx.Model(&quntity).Update("quantity", quntity.Quantity-orderbook[i].Quantity)
				if orm.Error != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Failed to update quantity",
					})
					return
				}
			}
		}

		//update
		ugorm := tx.Model(&order).Update("status", "waiting for confirmation")
		if ugorm.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed go to Finshed status",
			})
			return
		}
		tx.Commit()
	} else if body.Status == "waiting for cancelation" {
		if order.Status != "waiting for confirmation" && order.Status != "confirmed" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't change order status",
			})
			return
		}
		upgorm := initializers.DB.Model(&order).Update("status", "waiting for cancelation")
		if upgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to change order status",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong status",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "order status change success",
	})
}

func ChangeOrderStatus(c *gin.Context) {
	var body struct {
		Status string
		Id     uint
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	//find order
	var order models.Order
	initializers.DB.First(&order, "id = ? ", body.Id)
	if order.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "order not found",
		})
		return
	}

	if body.Status == "confirmed" || body.Status == "rejected" {
		if order.Status != "waiting for confirmation" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't change order status",
			})
			return
		}
		currentTime := time.Now()
		date :=  strconv.Itoa(int(currentTime.Year()))+"-"+strconv.Itoa(int(currentTime.Month()))+"-"+strconv.Itoa(int(currentTime.Day()))
		upgorm := initializers.DB.Model(&order).Updates(models.Order{
			Status: body.Status ,
			Confirm_date: date,})
		if upgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to change order status",
			})
			return
		}
	} else if body.Status == "canceled" {
		if order.Status != "waiting for cancelation" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Can't change order status",
			})
			return
		}
		upgorm := initializers.DB.Model(&order).Update("status", body.Status)
		if upgorm.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to change order status",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong status",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{
		"message": "order status change success",
	})
}








func CountSell(c *gin.Context) {
	//get params
	myurl, _ := url.Parse(c.Request.RequestURI)
	params, _ := url.ParseQuery(myurl.RawQuery)
	catregory_id, _ := strconv.ParseUint(params.Get("category_id"), 10, 64)
	start_date := params.Get("start_date")
	end_date := params.Get("end_date")

	var count struct {
		Count int
	}
	if catregory_id == 0 {
		initializers.DB.Model(&models.Order{}).Select("SUM(order_books.quantity) AS count").Where("status = 'confirmed' AND order.confirm_date >= ? AND order.confirm_date <= ? AND order_books.deleted_at IS NULL", start_date, end_date).Joins("JOIN order_books ON order_books.order_id = orders.id").Scan(&count)
	}else{
		initializers.DB.Model(&models.Order{}).Select("SUM(order_books.quantity) AS count").Where("status = 'confirmed' AND category_id = ? AND order.confirm_date >= ? AND order.confirm_date <= ? AND order_books.deleted_at IS NULL", catregory_id, start_date, end_date).Joins("JOIN order_books ON order_books.order_id = orders.id").Joins("JOIN quantity_books ON order_books.q_book_id = quantity_books.id").Joins("JOIN books ON quantity_books.book_id = books.id").Joins("JOIN category_books ON books.id = category_books.book_id").Scan(&count)		
	}
	c.JSON(http.StatusOK, count)
}

func PriceSell(c *gin.Context) {
	//get params
	myurl, _ := url.Parse(c.Request.RequestURI)
	params, _ := url.ParseQuery(myurl.RawQuery)
	catregory_id, _ := strconv.ParseUint(params.Get("category_id"), 10, 64)
	start_date := params.Get("start_date")
	end_date := params.Get("end_date")

	var price struct {
		Price int
	}

	if catregory_id == 0 {
		initializers.DB.Model(&models.Order{}).Select("SUM(order_books.quantity*quantity_books.s_price) AS price").Where("status = 'confirmed' AND order.confirm_date >= ? AND order.confirm_date <= ? AND order_books.deleted_at IS NULL", start_date, end_date).Joins("JOIN order_books ON order_books.order_id = orders.id").Joins("JOIN quantity_books ON order_books.q_book_id = quantity_books.id").Scan(&price)
	}else{
		initializers.DB.Model(&models.Order{}).Select("SUM(order_books.quantity*quantity_books.s_price) AS price").Where("status = 'confirmed' AND category_id = ? AND order.confirm_date >= ? AND order.confirm_date <= ? AND order_books.deleted_at IS NULL", catregory_id, start_date, end_date).Joins("JOIN order_books ON order_books.order_id = orders.id").Joins("JOIN quantity_books ON order_books.q_book_id = quantity_books.id").Joins("JOIN books ON quantity_books.book_id = books.id").Joins("JOIN category_books ON books.id = category_books.book_id").Scan(&price)		
	}
	c.JSON(http.StatusOK, price)
}

func GetRate(c *gin.Context) {
	//get params
	myurl, _ := url.Parse(c.Request.RequestURI)
	params, _ := url.ParseQuery(myurl.RawQuery)
	catregory_id, _ := strconv.ParseUint(params.Get("category_id"), 10, 64)
	start_date := params.Get("start_date")
	end_date := params.Get("end_date")

	var rate struct {
		Rate float64
	}

	if catregory_id != 0{
		initializers.DB.Table("rating_books").Select("AVG(rate) AS rate").Where("category_id = ? AND created_at >= ? AND created_at <= ? AND books.deleted_at IS NULL", catregory_id, start_date, end_date).Joins("JOIN books ON rating_books.book_id = books.id").Joins("JOIN books ON books.id = rating_books.book_id").Joins("JOIN category_books ON books.id = category_books.book_id").Scan(&rate)
	}else{
		initializers.DB.Table("rating_books").Select("AVG(rate) AS rate").Where("created_at >= ? AND created_at <= ?", start_date, end_date).Scan(&rate)
	}
	c.JSON(http.StatusOK, rate)
}	

func CountComment(c *gin.Context) {
	//get params
	myurl, _ := url.Parse(c.Request.RequestURI)
	params, _ := url.ParseQuery(myurl.RawQuery)
	start_date := params.Get("start_date")
	end_date := params.Get("end_date")

	var count struct {
		Count int
	}

	initializers.DB.Table("comments").Select("COUNT(id) AS count").Where("created_at >= ? AND created_at <= ?", start_date, end_date).Scan(&count)

	c.JSON(http.StatusOK, count)
}	

func Profit(c *gin.Context) {
	//get params
	myurl, _ := url.Parse(c.Request.RequestURI)
	params, _ := url.ParseQuery(myurl.RawQuery)
	start_date := params.Get("start_date")
	end_date := params.Get("end_date")

	var profit struct {
		Profit int
	}

	initializers.DB.Table("orders").Select("SUM(order_books.quantity*((CASE WHEN discount_books.id IS NOT NULL THEN (quantity_books.s_price*(100-discount_books.discount)/100) ELSE quantity_books.s_price END)-quantity_books.b_price)) AS profit").Where("status = 'confirmed' AND confirm_date >= ? AND confirm_date <= ?", start_date, end_date).Joins("JOIN order_books ON order_books.order_id = orders.id").Joins("JOIN quantity_books ON order_books.q_book_id = quantity_books.id").Joins("JOIN dicounts ON order_books.discount_id = discount_books.id").Scan(&profit)

	c.JSON(http.StatusOK, profit)
}