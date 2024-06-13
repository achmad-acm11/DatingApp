package middlewares

import (
	"DatingApp/helpers"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JwtTokenHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		var jwtKey = []byte(helpers.GetSecretKey())

		claims := jwt.MapClaims{}
		tokenString := helpers.GetTokenInHeader(context)

		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			context.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			context.Abort()
			return
		}
		location := helpers.GetLocalTime()
		timeNow := time.Now().In(location)

		expireTime, _ := claims["exp"].(float64)
		if int64(expireTime) < timeNow.Unix() {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			context.Abort()
			return
		}
		context.Set("userId", claims["data"].(map[string]interface{})["id"])
		context.Next()
	}
}
