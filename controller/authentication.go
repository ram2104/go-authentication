package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ram2104/go-authentication/initializer"
	model "github.com/ram2104/go-authentication/models"
	"golang.org/x/crypto/bcrypt"
)

/**
 */

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(ctx *gin.Context) {
	var req AuthRequest

	err := ctx.BindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": hashErr.Error(),
		})
		return
	}

	user := model.User{Email: req.Email, Password: string(hash)}

	result := initializer.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"count": result.RowsAffected,
	})
}

func Login(ctx *gin.Context) {
	var body AuthRequest

	err := ctx.BindJSON(&body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user model.User
	initializer.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "No account Found",
		})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fail to login user",
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{})
}
