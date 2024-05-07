package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/vova4o/go_final_project/internal/config"
)

type Password struct {
	Password string `json:"password"`
}

func SignIn(c *gin.Context) {
	var pwd Password
	if err := c.ShouldBindJSON(&pwd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if pwd.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "необходимо ввести пароль"})
		return
	}

	storedPasword := config.Password()

	if pwd.Password == storedPasword {

		token, err := CreateToken(pwd.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не верный пароль"})
		return
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// смотрим наличие пароля
		pass := config.Password()
		fmt.Println("Password:", pass) // Debug print
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := c.Request.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}
			fmt.Println("JWT:", jwt) // Debug print
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			valid, err = ValidateToken(jwt)
			fmt.Println("Valid:", valid) // Debug print
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating token"})
				c.Abort()
				return
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification required"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

var jwtKey = []byte("your_secret_key")

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.StandardClaims
}

func CreateToken(passwordHash string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		PasswordHash: passwordHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return false, err
	}

	if ok := token.Valid; ok {
		return true, nil
	}

	return false, nil
}
