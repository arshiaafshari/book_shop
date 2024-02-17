package main

import (
	"example/user/controllers"
	"example/user/initializers"
	"example/user/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
	//	initializers.ConnectToMinio()
}

func main() {
	r := gin.Default()

	//users*
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/logout", middleware.RequireAuth, controllers.Logout)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/getuser", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.GetUsers)
	r.DELETE("/deleteuser", middleware.RequireAuth, controllers.DeleteUser)
	r.PATCH("/updateuser", middleware.RequireAuth, controllers.UpdateUserInfo)
	r.PATCH("/setadmin", middleware.RequireAuth, middleware.IsMaster, controllers.SetAdmin)
	r.PATCH("/setmaster", middleware.RequireAuth, middleware.IsMaster, controllers.SetMaster)
	r.POST("/refresh", middleware.RequireAuthRefresh, controllers.GetNewTokens)

	//books*
	r.GET("/getbooks", controllers.GetBook)
	r.PUT("/addbook", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.AddBook)
	r.DELETE("/deletebook", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.DeleteBook)
	r.PATCH("/updatebook", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.UpdateBook)
	r.GET("/getbookinfo/:id", controllers.GetBookInfo)
	r.GET("/getbookcomments/:id", controllers.GetBookComment)
	r.GET("/getbookprice/:id", controllers.GetBookPrice)
	r.GET("/getbookcategory/:id", controllers.GetBookCategory)
	r.PUT("/addrating", middleware.RequireAuth, controllers.AddRating)

	//quantity*
	r.PUT("/addquantity", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.AddQuantity)
	r.PATCH("/updatequantity", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.UpdateQuantity)
	r.DELETE("/deletequantity", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.DeleteQuantity)
	r.GET("/getquantities", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.GetQuantity)

	//category*
	r.PUT("/addcategory", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.AddCategory)
	r.DELETE("/deletecategory", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.DeleteCategory)
	r.PATCH("/updatecategory", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.UpdateCategory)
	r.GET("/getcategories", controllers.GetCategory)

	//discount*
	r.PUT("/adddiscount", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.AddDiscount)
	r.PATCH("/updatediscount", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.UpdateDiscount)
	r.DELETE("/deletediscount", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.DeleteDiscount)
	r.GET("/getdiscounts", middleware.RequireAuth, middleware.IsPremium, controllers.GetDiscount)

	//comment*
	r.PUT("/addcomment", middleware.RequireAuth, controllers.AddComment)
	r.DELETE("/deletecomment", middleware.RequireAuth, controllers.DeleteComment)
	r.GET("/comments", middleware.RequireAuth, controllers.MyComments)
	r.GET("/checkcomment", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.CheckComment)
	r.PATCH("/confirmcomment", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.ConfirmComment)
	r.PUT("/likecomment", middleware.RequireAuth, controllers.LikeComment)

	//order*
	r.PUT("/setorder", middleware.RequireAuth, controllers.SetOrder)
	r.PUT("/addtoorder", middleware.RequireAuth, controllers.AddToOrder)
	r.DELETE("/deletefromorder", middleware.RequireAuth, controllers.DeleteFromOrder)
	r.DELETE("/deleteorder", middleware.RequireAuth, controllers.DeleteOrder)
	r.GET("/getwaitingorders", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.GetWaitingOrders)
	r.GET("/getorders", middleware.RequireAuth, controllers.GetOrders)
	r.GET("/getorderinfo/:id", middleware.RequireAuth, controllers.GetOrderInfo)
	r.PATCH("/userorderstatus", middleware.RequireAuth, controllers.UserChangeOrderStatus)
	r.PATCH("/orderstatus", middleware.RequireAuth, middleware.IsMasterOrAdmin, controllers.ChangeOrderStatus)

	//bookmark*
	r.PUT("/addbookmark", middleware.RequireAuth, controllers.AddBookmark)
	r.DELETE("/deletebookmark", middleware.RequireAuth, controllers.DeleteBookmark)
	r.GET("/getbookmarks", middleware.RequireAuth, controllers.GetBookmarks)
	
	//notification*
	r.GET("/newnotification", middleware.RequireAuth, controllers.NewNotifications)
	r.GET ("/getnotifications", middleware.RequireAuth, controllers.GetNotifications)

	//admin
	r.GET("/sellcount", controllers.CountSell)
	r.GET("/sellprice", controllers.PriceSell)
	r.GET("/commentcount", controllers.CountComment)
	r.GET("/getrating", controllers.GetRate)
	r.GET("/profit", controllers.Profit)


	r.Run()
}
