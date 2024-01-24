package controllers

import (
	"example/user/initializers"
	"example/user/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// get the information
	var body struct {
		UserName string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}
	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to hash password",
		})

		return
	}
	// create the user
	user := models.User{UserName: body.UserName, Password: string(hash),}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to create user",
		})

		return
	}
	//respond
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	//get the user and pss

	var body struct {
		UserName string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//look up req user
	var user models.User
	initializers.DB.First(&user, "user_name = ?", body.UserName)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username",
		})

		return
	}

	//check  pass

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid  password",
		})

		return
	}
	//Generate  jwt
	JwtId := NewJwtId(user.ID)
	//Access token
	Accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": JwtId,
		"exp": time.Now().Add(time.Second * 60).Unix(),
		"ref": false, 
	})

	AccesstokenString, err := Accesstoken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot create Access token",
		})

		return
	}

	//Refresh token
	Refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": JwtId,
		"exp": time.Now().Add(time.Hour).Unix(),
		"ref": true,
	})

	RefreshtokenString, err := Refreshtoken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot create Refresh token",
		})
	}
	//put JwtId in database
	qgorm := initializers.DB.Model(&user).Update("jwt_id", JwtId)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot put JwtId in database",
		})
		return
	}
	//send it back

	c.JSON(http.StatusOK, gin.H{
		"Access token": AccesstokenString,
		"Refresh token": RefreshtokenString,
	})
}

func Logout(c *gin.Context) {
	//find the user
	var user models.User
	qgorm := initializers.DB.First(&user, "id = ?", c.MustGet("USERID"))
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to find user",
		})
		return
	}
	//change the jwt id
	upgorm := initializers.DB.Model(&user).Update("jwt_id", NewJwtId(user.ID))
	if upgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to logout",
		})
		return
	}
	//send it back
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out",
	})
}

func UpdateUserInfo(c *gin.Context) {
	//get info from body
	var body struct {
		ID       uint
		UserName string
		Email    string
		FName    string
		LName    string 
		City     string 
		Birth    string 
		Gender   string 
		Phone    string 
		Address  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	return	
	}
	//find user
	var user models.User
	initializers.DB.First(&user, "id = ?", body.ID)

	if user.ID == 0 {
	    c.JSON(http.StatusBadRequest, gin.H{
	        "error": "Failed to find user",
	    })
	    return
	}
	//check permision
	role, _ := c.Get("Role")
	userid, _ := c.Get("USERID")

	if user.Role == "user" {

		if role != "master" && role != "admin" && userid != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Don't have permision to update this user",
			})

			return
		}
	} else if user.Role == "admin" {

		if role != "master" && userid != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Don't have permision to update this admin",
			})

			return
		}
	}else if user.Role == "master" {
		if userid != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Don't have permision to update this master",
			})

			return
		}
	}	
	//check username not exist
	if body.UserName != "" {
		var usercheck models.User
		initializers.DB.First(&usercheck, "user_name = ?", body.UserName)

		if usercheck.ID != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username already exist",
			})
		}
		return
	}

	//update
	upgorm := initializers.DB.Model(&user).Updates(models.User{
		UserName: body.UserName,
		Email:    body.Email,
		FName:    body.FName,
		LName:    body.LName,
		City:     body.City,
		Birth:    body.Birth,
		Gender:   body.Gender,
		Phone:    body.Phone,
		Address:  body.Address,
	})

	if upgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
	})
}

func DeleteUser(c *gin.Context) {

	var body struct {
		Id string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	id, _ := strconv.Atoi(body.Id)
	//look up req user
	var user models.User
	initializers.DB.First(&user, "id = ?", id)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id",
		})

		return
	}

	//check permision
	role, _ := c.Get("Role")
	userid, _ := c.Get("USERID")

	if user.Role == "user" {

		if role != "master" && role != "admin" && userid != user.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Don't have permision to delete this user",
			})

			return
		}
	} else if user.Role == "admin" || user.Role == "master" {

		if role != "master" && userid != user.ID {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Don't have permision to delete this master or admin",
			})

			return
		}
	}

	// soft delete
	qgorm := initializers.DB.Delete(&user)

	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete user",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user delete success",
	})
}

func SetAdmin(c *gin.Context) {
	//get user name
	var body struct {
		ID uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}

	//look up req user
	var user models.User
	initializers.DB.First(&user, "id = ?", body.ID)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username",
		})

		return
	}

	// set admin
	qgorm := initializers.DB.Model(&user).Update("role", "admin")
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to set admin",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Set admin success",
	})
}

func SetMaster(c *gin.Context) {
	//get user name
	var body struct {
		ID uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//look up req user
	var user models.User
	initializers.DB.First(&user, "id = ?", body.ID)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username",
		})

		return
	}

	// set master
	qgorm := initializers.DB.Model(&user).Update("role", "master")
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to set master",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Set master success",
	})
}

func GetNewTokens(c *gin.Context) {

	//get ID
	UserId:= c.GetUint("USERID")

	//Generate new jwt
	JwtId := NewJwtId(UserId)
	//Access token
	Accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": JwtId,
		"exp": time.Now().Add(time.Second * 60).Unix(),
		"ref": false, 
	})

	AccesstokenString, err := Accesstoken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot create Access token",
		})

		return
	}

	//Refresh token
	Refreshtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": JwtId,
		"exp": time.Now().Add(time.Hour).Unix(),
		"ref": true,
	})

	RefreshtokenString, err := Refreshtoken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot create Refresh token",
		})
	}
	//find user
	var user models.User
	initializers.DB.First(&user, "id = ?", UserId)
	//chandge jwt id
	initializers.DB.Model(&user).Update("jwt_id", JwtId)

	//send it back

	c.JSON(http.StatusOK, gin.H{
		"Access token": AccesstokenString,
		"Refresh token": RefreshtokenString,
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("USERID")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func GetUsers(c *gin.Context) {

	var users []models.User
	qgorm := initializers.DB.Find(&users)
	if qgorm.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get users",
		})
		return
	}
	
	var users2 []models.Users
	for i := 0; i < len(users); i++ {
		users2 = append(users2, models.Users{
			ID:        users[i].ID,
			UserName:  users[i].UserName,
			Role:      users[i].Role,
			Email:     users[i].Email,
			FName:     users[i].FName,
			LName:     users[i].LName,
			City:      users[i].City,
			Birth:     users[i].Birth,
			Gender:    users[i].Gender,
			Phone:     users[i].Phone,
			Address:   users[i].Address,
			Score:     users[i].Score,
		})
	}
	c.JSON(http.StatusOK, users2)
}







func NewJwtId (id uint) (string) {
	jwtid := strconv.FormatUint(uint64(id), 10) +"."+ (strconv.FormatInt(time.Now().Unix(), 10))[6:10]
	return jwtid
}