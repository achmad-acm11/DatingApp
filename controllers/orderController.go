package controllers

import (
	constants "DatingApp/constansts"
	"DatingApp/helpers"
	"DatingApp/requests"
	"DatingApp/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type orderController struct {
	service services.OrderService
}

func NewOrderController(service services.OrderService) *orderController {
	return &orderController{service: service}
}

func (o orderController) CreateHandler(ctx *gin.Context) {
	var reqOrder requests.CreateOrderRequest
	err := ctx.ShouldBindJSON(&reqOrder)
	helpers.ErrorHandler(err)

	orderResponse := o.service.CreateOrder(ctx, reqOrder)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success create order", "success", http.StatusOK, orderResponse))
}

func (o orderController) GetDetailHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helpers.ErrorHandler(err)

	orderResponse := o.service.GetOneOrder(ctx, id)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get order", "success", http.StatusOK, orderResponse))
}

func (p orderController) DeleteOneHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helpers.ErrorHandlerValidator(err)

	p.service.DeleteOneOrder(ctx, id)

	message := make(map[string]string)
	message["message"] = constants.OrderDeleted

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success delete order", "success", http.StatusOK, message))
}
