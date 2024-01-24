package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"example/user/initializers"
	"example/user/models"
	"example/user/controllers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	// get the token from header
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	// check if token is empty

	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "No token provided",
		})
		return
	}

	//parse takes the token string and a function for looking up the key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error" : err.Error(),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check not refresh
		if claims["ref"].(bool) {
		    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		        "error": "refresh token not allowed",
		    })
		    return
		}



		// find the user with token sub
		var user models.User
		initializers.DB.First(&user, "jwt_id = ?", claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}
		// attach to request
		c.Set("Role", user.Role)
		c.Set("USERID", user.ID)
		c.Set("Score", user.Score)

		// continue
		c.Next()

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
}

func RequireAuthRefresh(c *gin.Context) {
	//Get the refresh token
	var body struct {
		RefreshToken string
	}

	if c.Bind(&body) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//parse takes the token string and a function for looking up the key
	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// check not refresh
		if !claims["ref"].(bool) {
		    c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		        "error": "Access token not allowed",
		    })
		    return
		}

		// find the user with token sub
		var user models.User
		initializers.DB.First(&user, "jwt_id = ?", claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token blocked. try to block all token",
			})
			JwtId  := claims["sub"].(string)
			UserId := (strings.Split(JwtId, "."))[0]
			fgorm := initializers.DB.First(&user, "id = ?", UserId)
			if fgorm.Error != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": ": user not found",
				})
				return
			}
			upgorm := initializers.DB.Model(&user).Update("jwt_id", controllers.NewJwtId(user.ID))
			if upgorm.Error != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "Failed to block all token",
				})
				return
			}
			return
		}

		// attach to request
		c.Set("USERID", user.ID)

		// continue
		c.Next()

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
}

func IsMasterOrAdmin(c *gin.Context) {
	role , _ := c.Get("Role")
	if role != "master" && role != "admin" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Do not have permission",
		})
		return
	}
	c.Next()
}	

func IsMaster(c *gin.Context) {
	role, _ := c.Get("Role")
	if role != "master" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Do not have permission",
		})
		return
	}
	c.Next()
}

func IsPremium(c *gin.Context) {
	score := c.MustGet("Score").(int)
	if score < 100 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "You are not premium user",
		})
		return
	}
	c.Next()
}