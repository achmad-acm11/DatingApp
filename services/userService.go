package services

import (
	"DatingApp/constansts"
	"DatingApp/entities"
	"DatingApp/exceptions"
	"DatingApp/helpers"
	"DatingApp/repositories"
	"DatingApp/requests"
	"DatingApp/responses"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type UserService interface {
	GetAllUser(ctx *gin.Context) []entities.User
	SignUpUser(ctx *gin.Context, request requests.SignUpUserRequest) entities.User
	GetOneUser(ctx *gin.Context, userId int) entities.User
	GetCouple(ctx *gin.Context, userId int) entities.User
	LikeCouple(ctx *gin.Context, userId int) entities.User
	PassCouple(ctx *gin.Context, userId int) entities.User
	//CheckOTP(ctx *gin.Context, request requests.CheckOTPRequest)
	//CheckUser(ctx *gin.Context, request requests.CheckUserRequest) map[string]string
	//SendOTP(ctx *gin.Context, request requests.SendOTPUserRequest) map[string]string
	Login(ctx *gin.Context, request requests.LoginRequest) responses.TokenResponse
	//UpdateUser(ctx *gin.Context, request userRequestsParam.UpdateUserRequest, userId int) responses.UserResponse
	DeleteOneUser(ctx *gin.Context, id int)
}

type userService struct {
	repo         repositories.UserRepository
	repoUserFind repositories.UserFindRepository
	db           *gorm.DB
	validate     *validator.Validate
}

func NewUserService(
	db *gorm.DB,
	validate *validator.Validate,
	repo repositories.UserRepository,
	repoUserFind repositories.UserFindRepository,
) userService {
	return userService{
		repo:         repo,
		repoUserFind: repoUserFind,
		db:           db,
		validate:     validate,
	}
}

func (u userService) GetAllUser(ctx *gin.Context) []entities.User {
	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	users := u.repo.GetAll(ctx, tx)

	return users
}

func (u userService) SignUpUser(ctx *gin.Context, request requests.SignUpUserRequest) entities.User {
	err := u.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	userEmail := u.repo.GetOneByEmail(ctx, tx, request.Email)
	if userEmail.Email != "" {
		panic(exceptions.NewConflictError(errors.New(constants.UserConflict).Error()))
	}

	password := u.generatePassword(request.Password)

	user := u.repo.Create(ctx, tx, entities.User{
		Name:        request.Name,
		Gender:      request.Gender,
		PhoneNumber: strings.TrimSpace(request.PhoneNumber),
		Email:       strings.TrimSpace(request.Email),
		Password:    string(password),
	})

	return user
}

func (u userService) Login(ctx *gin.Context, request requests.LoginRequest) responses.TokenResponse {
	err := u.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	userData := u.repo.GetOneByEmail(ctx, tx, request.Email)
	u.checkUserExistWithPanic(userData, exceptions.NewUnauthorizedError(errors.New(constants.LoginFailed).Error()))

	u.checkPassword(request.Password, userData.Password, exceptions.NewUnauthorizedError(errors.New(constants.LoginFailed).Error()))

	token := u.generateToken(userData)

	return responses.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   fmt.Sprintf("%d hours", helpers.GetExpiredNum()),
		User:        userData,
	}
}

func (u userService) GetOneUser(ctx *gin.Context, userId int) entities.User {
	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	return user
}

func (u userService) DeleteOneUser(ctx *gin.Context, id int) {
	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, id)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	u.repo.DeleteOne(ctx, tx, user)
}

func (u userService) GetCouple(ctx *gin.Context, userId int) entities.User {
	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	userFind := u.repoUserFind.GetOneByUserId(ctx, tx, user.Id)

	timeNow := helpers.GetLocalDateNow()
	if userFind.Id == 0 {
		userFind = u.repoUserFind.Create(ctx, tx, user)
	}

	if userFind.Date.Unix() < timeNow.Unix() {
		userFind = u.repoUserFind.Reset(ctx, tx, userFind, user)
	}

	if user.IsPremium == false && userFind.CurrentState >= 10 {
		panic(exceptions.NewPaymentRequiredError(errors.New(constants.PaymentRequired).Error()))
	}
	userQueue := strings.Split(userFind.UserQueue, "|")
	userIdCouple, _ := strconv.Atoi(userQueue[userFind.CurrentState-1])

	userCouple := u.repo.GetOneById(ctx, tx, userIdCouple)

	return userCouple
}

func (u userService) LikeCouple(ctx *gin.Context, userId int) entities.User {
	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	timeNow := helpers.GetLocalDateNow()
	userFind := u.repoUserFind.GetOneCurrentFind(ctx, tx, user.Id, timeNow)

	if userFind.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserCoupleNotFound).Error()))
	}

	if user.IsPremium == false && userFind.CurrentState >= 10 {
		panic(exceptions.NewPaymentRequiredError(errors.New(constants.PaymentRequired).Error()))
	}
	userQueue := strings.Split(userFind.UserQueue, "|")
	userIdCouple, _ := strconv.Atoi(userQueue[(userFind.CurrentState-1)+1])

	userCouple := u.repo.GetOneById(ctx, tx, userIdCouple)

	userFind.LikeCount += 1
	userFind.CurrentState += 1
	u.repoUserFind.Update(ctx, tx, userFind)

	return userCouple
}

func (u userService) PassCouple(ctx *gin.Context, userId int) entities.User {
	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	timeNow := helpers.GetLocalDateNow()
	userFind := u.repoUserFind.GetOneCurrentFind(ctx, tx, user.Id, timeNow)

	if userFind.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserCoupleNotFound).Error()))
	}

	if user.IsPremium == false && userFind.CurrentState >= 10 {
		panic(exceptions.NewPaymentRequiredError(errors.New(constants.PaymentRequired).Error()))
	}
	userQueue := strings.Split(userFind.UserQueue, "|")
	userIdCouple, _ := strconv.Atoi(userQueue[(userFind.CurrentState-1)+1])

	userCouple := u.repo.GetOneById(ctx, tx, userIdCouple)

	userFind.PassCount += 1
	userFind.CurrentState += 1
	u.repoUserFind.Update(ctx, tx, userFind)

	return userCouple
}

func (u userService) generatePassword(rawPassword string) []byte {
	password, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 14)
	helpers.ErrorHandler(err)
	return password
}

func (u userService) checkUserExistWithPanic(user entities.User, messageError interface{}) {
	if user.Id == 0 {
		panic(messageError)
	}
}

func (u userService) checkPassword(requestPassword string, userPassword string, errorException interface{}) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(requestPassword))
	if err != nil {
		panic(errorException)
	}
}

func (u userService) generateToken(data interface{}) string {
	token, err := helpers.GenerateToken(map[string]interface{}{
		"data": data,
	})
	helpers.ErrorHandler(err)
	return token
}

//func (u userService) UpdateUser(ctx *gin.Context, request userRequestsParam.UpdateUserRequest, userId int) responses.UserResponse {
//	u.stdLog.NameFunc = "UpdateUser"
//	u.stdLog.StartFunction(map[string]interface{}{"request": request, "userId": userId})
//
//	err := u.validate.Struct(request)
//	helper.ErrorHandlerValidator(err)
//
//	tx := u.db.Begin()
//	defer helper.CommitOrRollback(tx)
//
//	userOld := u.repo.GetOneById(ctx, tx, userId)
//	u.checkUserExistWithPanic(userOld, exception.NewNotFoundError(errors.New(constant.UserNotFound).Error()))
//
//	paramBuilder := u.paramBuilder.CreateUserUpdateRequestParam().InitUserOld(userOld).
//		AddEmail(request.Email, u.validate).
//		AddPhoneNumber(request.PhoneNumber, u.validate).
//		AddPassword(request.Password)
//
//	//// Verification format date
//	//if request.VerificationDate != "" {
//	//	date, _ := time.Parse("2006-01-02", request.VerificationDate)
//	//	userOld.VerificationDate = date
//	//}
//	//userOld.Name = request.Name
//	////userOld.Role = request.Role
//	//userOld.PhoneNumber = request.PhoneNumber
//	//userOld.Email = request.Email
//
//	user := u.repoRedis.Update(ctx, tx, paramBuilder.ResultRequest())
//
//	go u.repoStudent.UpdateByUserId(requests.StudentUpdateRequest{
//		Email:       user.Email,
//		PhoneNumber: user.PhoneNumber,
//	}, user.Id, helper.GetTokenInHeader(ctx))
//
//	// End log
//	u.stdLog.NameFunc = "UpdateUser"
//	u.stdLog.EndFunction(user)
//
//	response := responses.NewUserResponseBuilder()
//	return response.Default(user).Result()
//}
