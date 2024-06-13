package helpers

import (
	"github.com/gin-gonic/gin"
)

func GetUserIdContext(ctx *gin.Context) int {
	userIdAny := ctx.MustGet("userId").(float64)
	return int(userIdAny)
}
