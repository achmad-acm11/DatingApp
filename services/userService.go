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
	GetMyUser(ctx *gin.Context) entities.User
	GetAllMatch(ctx *gin.Context) entities.User
	LikeMatch(ctx *gin.Context) entities.User
	PassMatch(ctx *gin.Context) entities.User
	Login(ctx *gin.Context, request requests.LoginRequest) responses.TokenResponse
	DeleteOneUser(ctx *gin.Context, id int)
}

type userService struct {
	repo           repositories.UserRepository
	repoMatchQueue repositories.MatchQueueRepository
	db             *gorm.DB
	validate       *validator.Validate
	stdLog         *helpers.StandartLog
}

func NewUserService(
	db *gorm.DB,
	validate *validator.Validate,
	repo repositories.UserRepository,
	repoMatchQueue repositories.MatchQueueRepository,
) userService {
	return userService{
		repo:           repo,
		repoMatchQueue: repoMatchQueue,
		db:             db,
		validate:       validate,
		stdLog:         helpers.NewStandardLog(constants.User, constants.Service),
	}
}

func (u userService) GetAllUser(ctx *gin.Context) []entities.User {
	u.stdLog.NameFunc = "GetAllUser"
	u.stdLog.StartFunction(nil)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	users := u.repo.GetAll(ctx, tx)

	u.stdLog.NameFunc = "GetAllUser"
	u.stdLog.EndFunction(nil)

	return users
}

func (u userService) SignUpUser(ctx *gin.Context, request requests.SignUpUserRequest) entities.User {
	u.stdLog.NameFunc = "SignUpUser"
	u.stdLog.StartFunction(request)

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

	u.stdLog.NameFunc = "SignUpUser"
	u.stdLog.EndFunction(user)

	return user
}

func (u userService) Login(ctx *gin.Context, request requests.LoginRequest) responses.TokenResponse {
	u.stdLog.NameFunc = "Login"
	u.stdLog.StartFunction(request)

	err := u.validate.Struct(request)
	helpers.ErrorHandlerValidator(err)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	userData := u.repo.GetOneByEmail(ctx, tx, request.Email)
	u.checkUserExistWithPanic(userData, exceptions.NewUnauthorizedError(errors.New(constants.LoginFailed).Error()))

	u.checkPassword(request.Password, userData.Password, exceptions.NewUnauthorizedError(errors.New(constants.LoginFailed).Error()))

	token := u.generateToken(userData)

	u.stdLog.NameFunc = "Login"
	u.stdLog.EndFunction(userData)

	return responses.TokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   fmt.Sprintf("%d hours", helpers.GetExpiredNum()),
		User:        userData,
	}
}

func (u userService) GetMyUser(ctx *gin.Context) entities.User {
	userId := helpers.GetUserIdContext(ctx)

	u.stdLog.NameFunc = "GetMyUser"
	u.stdLog.StartFunction(userId)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	u.stdLog.NameFunc = "GetMyUser"
	u.stdLog.EndFunction(user)

	return user
}

func (u userService) GetOneUser(ctx *gin.Context, userId int) entities.User {
	u.stdLog.NameFunc = "GetOneUser"
	u.stdLog.StartFunction(userId)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	u.stdLog.NameFunc = "GetOneUser"
	u.stdLog.EndFunction(user)

	return user
}

func (u userService) DeleteOneUser(ctx *gin.Context, id int) {
	u.stdLog.NameFunc = "DeleteOneUser"
	u.stdLog.StartFunction(id)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, id)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	u.repo.DeleteOne(ctx, tx, user)

	u.stdLog.NameFunc = "DeleteOneUser"
	u.stdLog.EndFunction(nil)
}

func (u userService) GetAllMatch(ctx *gin.Context) entities.User {
	userId := helpers.GetUserIdContext(ctx)

	u.stdLog.NameFunc = "GetAllMatch"
	u.stdLog.StartFunction(userId)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	matchQueue := u.repoMatchQueue.GetOneByUserId(ctx, tx, user.Id)

	timeNow := helpers.GetLocalDateNow()
	if matchQueue.Id == 0 {
		matchQueue = u.repoMatchQueue.Create(ctx, tx, user)
	}

	if matchQueue.Date.Unix() < timeNow.Unix() {
		matchQueue = u.repoMatchQueue.Reset(ctx, tx, matchQueue, user)
	}

	if user.IsPremium == false && matchQueue.CurrentState >= 10 {
		panic(exceptions.NewPaymentRequiredError(errors.New(constants.PaymentRequired).Error()))
	} else if user.IsPremium == true && matchQueue.CurrentState >= 10 && matchQueue.CurrentState%10 == 0 {
		matchQueue = u.repoMatchQueue.AppendQueue(ctx, tx, matchQueue, user)
	}

	if (matchQueue.CurrentState-1)+1 >= len(matchQueue.UserQueue) {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserMatchNoMore).Error()))
	}

	userQueue := strings.Split(matchQueue.UserQueue, "|")
	if (matchQueue.CurrentState-1)+1 >= len(userQueue) {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserMatchNotFound).Error()))
	}
	userIdMatch, _ := strconv.Atoi(userQueue[matchQueue.CurrentState-1])

	userMatch := u.repo.GetOneById(ctx, tx, userIdMatch)

	u.stdLog.NameFunc = "GetAllMatch"
	u.stdLog.EndFunction(userMatch)

	return userMatch
}

func (u userService) LikeMatch(ctx *gin.Context) entities.User {
	userId := helpers.GetUserIdContext(ctx)

	u.stdLog.NameFunc = "LikeMatch"
	u.stdLog.StartFunction(userId)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	timeNow := helpers.GetLocalDateNow()
	matchQueue := u.repoMatchQueue.GetOneCurrentQueue(ctx, tx, user.Id, timeNow)

	if matchQueue.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserMatchNotFound).Error()))
	}

	matchQueue.LikeCount += 1
	u.repoMatchQueue.Update(ctx, tx, matchQueue)

	if user.IsPremium == false && matchQueue.CurrentState >= 10 {
		panic(exceptions.NewPaymentRequiredError(errors.New(constants.PaymentRequired).Error()))
	} else if user.IsPremium == true && matchQueue.CurrentState >= 10 && matchQueue.CurrentState%10 == 0 {
		matchQueue = u.repoMatchQueue.AppendQueue(ctx, tx, matchQueue, user)
	}

	userQueue := strings.Split(matchQueue.UserQueue, "|")
	if (matchQueue.CurrentState-1)+1 >= len(userQueue) {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserMatchNoMore).Error()))
	}
	userIdMatch, _ := strconv.Atoi(userQueue[(matchQueue.CurrentState-1)+1])

	userMatch := u.repo.GetOneById(ctx, tx, userIdMatch)

	matchQueue.CurrentState += 1
	u.repoMatchQueue.Update(ctx, tx, matchQueue)

	u.stdLog.NameFunc = "LikeMatch"
	u.stdLog.EndFunction(userMatch)

	return userMatch
}

func (u userService) PassMatch(ctx *gin.Context) entities.User {
	userId := helpers.GetUserIdContext(ctx)

	u.stdLog.NameFunc = "PassMatch"
	u.stdLog.StartFunction(userId)

	tx := u.db.Begin()
	defer helpers.CommitOrRollback(tx)

	user := u.repo.GetOneById(ctx, tx, userId)
	u.checkUserExistWithPanic(user, exceptions.NewNotFoundError(errors.New(constants.UserNotFound).Error()))

	timeNow := helpers.GetLocalDateNow()
	matchQueue := u.repoMatchQueue.GetOneCurrentQueue(ctx, tx, user.Id, timeNow)

	if matchQueue.Id == 0 {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserMatchNotFound).Error()))
	}
	matchQueue.PassCount += 1
	u.repoMatchQueue.Update(ctx, tx, matchQueue)

	if user.IsPremium == false && matchQueue.CurrentState >= 10 {
		panic(exceptions.NewPaymentRequiredError(errors.New(constants.PaymentRequired).Error()))
	} else if user.IsPremium == true && matchQueue.CurrentState >= 10 && matchQueue.CurrentState%10 == 0 {
		matchQueue = u.repoMatchQueue.AppendQueue(ctx, tx, matchQueue, user)
	}

	userQueue := strings.Split(matchQueue.UserQueue, "|")
	if (matchQueue.CurrentState-1)+1 >= len(userQueue) {
		panic(exceptions.NewNotFoundError(errors.New(constants.UserMatchNoMore).Error()))
	}
	userIdMatch, _ := strconv.Atoi(userQueue[(matchQueue.CurrentState-1)+1])

	userMatch := u.repo.GetOneById(ctx, tx, userIdMatch)

	matchQueue.CurrentState += 1
	u.repoMatchQueue.Update(ctx, tx, matchQueue)

	u.stdLog.NameFunc = "PassMatch"
	u.stdLog.EndFunction(userMatch)

	return userMatch
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
