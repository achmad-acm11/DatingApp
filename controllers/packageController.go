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

type packageController struct {
	service services.PackageService
}

func NewPackageController(service services.PackageService) *packageController {
	return &packageController{service: service}
}

func (p packageController) GetAllHandler(ctx *gin.Context) {
	packagesResponse := p.service.GetAllPackage(ctx)
	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get packages", "success", http.StatusOK, packagesResponse))
}

func (p packageController) CreateHandler(ctx *gin.Context) {
	var reqPackage requests.CreatePackageRequest
	err := ctx.ShouldBindJSON(&reqPackage)
	helpers.ErrorHandler(err)

	packageResponse := p.service.CreatePackage(ctx, reqPackage)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success create package", "success", http.StatusOK, packageResponse))
}

func (p packageController) GetDetailHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helpers.ErrorHandler(err)

	packageResponse := p.service.GetOnePackage(ctx, id)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success get package", "success", http.StatusOK, packageResponse))
}

func (p packageController) UpdateHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helpers.ErrorHandler(err)

	var reqPackage requests.UpdatePackageRequest
	err = ctx.ShouldBindJSON(&reqPackage)
	helpers.ErrorHandler(err)

	packageResponse := p.service.UpdatePackage(ctx, reqPackage, id)

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success update package", "success", http.StatusOK, packageResponse))
}

func (p packageController) DeleteOneHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	helpers.ErrorHandlerValidator(err)

	p.service.DeleteOnePackage(ctx, id)

	message := make(map[string]string)
	message["message"] = constants.PackageDeleted

	ctx.JSON(http.StatusOK, helpers.APIResponse("Success delete package", "success", http.StatusOK, message))
}
