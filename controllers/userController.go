package controllers

import (
	constants "DatingApp/constansts"
	"DatingApp/helpers"
	"DatingApp/requests"
	"DatingApp/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *userController {
	return &userController{service: service}
}

func (c userController) GetAllHandler(ctx *gin.Context) {
	usersResponse := c.service.GetAllUser(ctx)
	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get users", "success", http.StatusOK, usersResponse))
}

func (c userController) SignUpHandler(ctx *gin.Context) {
	var reqUser requests.SignUpUserRequest
	err := ctx.ShouldBindJSON(&reqUser)
	helpers.ErrorHandler(err)

	userResponse := c.service.SignUpUser(ctx, reqUser)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success create user", "success", http.StatusOK, userResponse))
}

func (c userController) LoginHandler(ctx *gin.Context) {
	var reqUser requests.LoginRequest
	err := ctx.ShouldBindJSON(&reqUser)
	helpers.ErrorHandler(err)

	tokenResponse := c.service.Login(ctx, reqUser)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success login user", "success", http.StatusOK, tokenResponse))
}

func (c userController) GetMyProfileHandler(ctx *gin.Context) {
	userResponse := c.service.GetMyUser(ctx)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get my user profile", "success", http.StatusOK, userResponse))
}

func (c userController) GetProfileHandler(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	helpers.ErrorHandler(err)

	userResponse := c.service.GetOneUser(ctx, userId)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get user profile", "success", http.StatusOK, userResponse))
}

func (c userController) GetMatchesHandler(ctx *gin.Context) {
	userResponse := c.service.GetAllMatch(ctx)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get matches", "success", http.StatusOK, userResponse))
}

func (c userController) LikeMatchHandler(ctx *gin.Context) {
	userResponse := c.service.LikeMatch(ctx)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success like match", "success", http.StatusOK, userResponse))
}

func (c userController) PassMatchHandler(ctx *gin.Context) {
	userResponse := c.service.PassMatch(ctx)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success pass match", "success", http.StatusOK, userResponse))
}

func (c userController) DeleteOneHandler(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	helpers.ErrorHandlerValidator(err)

	c.service.DeleteOneUser(ctx, userId)

	message := make(map[string]string)
	message["message"] = constants.UserDeleted

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success delete user", "success", http.StatusOK, message))
}
