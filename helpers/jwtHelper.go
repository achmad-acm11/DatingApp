package helpers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

func GetSecretKey() string {
	if gin.Mode() != gin.TestMode {
		if os.Getenv("APP_ENV") == "" {
			errEnv := godotenv.Load(".env")

			ErrorHandler(errEnv)
		}
	}

	return os.Getenv("JWT_SECRET")
}

func getIssKey() string {
	if gin.Mode() != gin.TestMode {
		if os.Getenv("APP_ENV") == "" {
			errEnv := godotenv.Load(".env")

			ErrorHandler(errEnv)
		}
	}

	return os.Getenv("JWT_KEY_ISS")
}

func GetExpiredNum() int {
	if gin.Mode() != gin.TestMode {
		if os.Getenv("APP_ENV") == "" {
			errEnv := godotenv.Load(".env")

			ErrorHandler(errEnv)
		}
	}
	expiredTime, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED"))
	return expiredTime
}

func getExpired() time.Duration {
	expiredTime := GetExpiredNum()
	return time.Duration(expiredTime)
}

func GenerateToken(data map[string]interface{}) (string, error) {
	claim := jwt.MapClaims{}
	timeData := time.Now()
	data["iss"] = getIssKey()
	data["iat"] = timeData.Unix()
	data["nbf"] = timeData.Unix()
	data["exp"] = timeData.Add(getExpired() * time.Hour).Unix()
	claim = data
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(GetSecretKey()))

	if err != nil {
		return token, err
	}
	return token, nil
}

func GetTokenInHeader(ctx *gin.Context) string {
	if gin.Mode() != gin.TestMode {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token := authHeader[7:]
				return token
			}
		}
	}
	return ""
}
