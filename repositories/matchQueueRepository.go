package repositories

import (
	"DatingApp/entities"
	"DatingApp/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type MatchQueueRepository interface {
	GetAll(ctx *gin.Context, db *gorm.DB) []entities.MatchQueue
	GetOneByUserId(ctx *gin.Context, db *gorm.DB, userId int) entities.MatchQueue
	GetOneCurrentQueue(ctx *gin.Context, db *gorm.DB, userId int, date time.Time) entities.MatchQueue
	Create(ctx *gin.Context, db *gorm.DB, user entities.User) entities.MatchQueue
	Update(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue) entities.MatchQueue
	Reset(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue, user entities.User) entities.MatchQueue
	AppendQueue(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue, user entities.User) entities.MatchQueue
	DeleteOne(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue)
}

type queueRepository struct {
	repoUser UserRepository
}

func NewMatchQueueRepository(repoUser UserRepository) queueRepository {
	return queueRepository{
		repoUser: repoUser,
	}
}

func (f queueRepository) GetAll(ctx *gin.Context, db *gorm.DB) []entities.MatchQueue {
	queues := []entities.MatchQueue{}
	tx := db.WithContext(ctx)

	err := tx.Limit(500).Find(&queues).Error
	helpers.ErrorHandler(err)

	return queues
}

func (f queueRepository) GetOneByUserId(ctx *gin.Context, db *gorm.DB, userId int) entities.MatchQueue {
	queue := entities.MatchQueue{}

	err := db.WithContext(ctx).Where("user_id = ?", userId).Find(&queue).Error
	helpers.ErrorHandler(err)

	return queue
}

func (f queueRepository) GetOneCurrentQueue(ctx *gin.Context, db *gorm.DB, userId int, date time.Time) entities.MatchQueue {
	queue := entities.MatchQueue{}

	err := db.WithContext(ctx).Where("user_id = ?", userId).Where("date = ?", date).Find(&queue).Error
	helpers.ErrorHandler(err)

	return queue
}

func (f queueRepository) Create(ctx *gin.Context, db *gorm.DB, user entities.User) entities.MatchQueue {
	timeNow := helpers.GetLocalDateNow()

	userQueue := f.getUserQueue(ctx, db, user.Id, user.Gender, []string{})

	queue := entities.MatchQueue{
		UserId:       user.Id,
		PassCount:    0,
		LikeCount:    0,
		CurrentState: 1,
		UserQueue:    strings.Join(userQueue, "|"),
		Date:         timeNow,
	}

	err := db.WithContext(ctx).Create(&queue).Error
	helpers.ErrorHandler(err)

	return queue
}

func (f queueRepository) Update(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue) entities.MatchQueue {
	err := db.WithContext(ctx).Save(&queue).Error
	helpers.ErrorHandler(err)

	return queue
}

func (f queueRepository) Reset(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue, user entities.User) entities.MatchQueue {
	timeNow := helpers.GetLocalDateNow()

	userQueue := f.getUserQueue(ctx, db, user.Id, user.Gender, []string{})

	queue.CurrentState = 1
	queue.PassCount = 0
	queue.LikeCount = 0
	queue.UserQueue = strings.Join(userQueue, "|")
	queue.Date = timeNow

	queue = f.Update(ctx, db, queue)

	return queue
}

func (f queueRepository) AppendQueue(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue, user entities.User) entities.MatchQueue {
	timeNow := helpers.GetLocalDateNow()

	userQueueOld := strings.Split(queue.UserQueue, "|")
	userQueue := f.getUserQueue(ctx, db, user.Id, user.Gender, userQueueOld)

	appendQueue := append(userQueueOld, userQueue...)
	queue.UserQueue = strings.Join(appendQueue, "|")
	queue.Date = timeNow

	queue = f.Update(ctx, db, queue)

	return queue
}

func (f queueRepository) DeleteOne(ctx *gin.Context, db *gorm.DB, queue entities.MatchQueue) {
	err := db.WithContext(ctx).Delete(&queue).Error
	helpers.ErrorHandler(err)
}

func (f queueRepository) getUserQueue(ctx *gin.Context, db *gorm.DB, userId int, gender string, userIds []string) []string {
	usersMatch := f.repoUser.GetAllTenMatch(ctx, db, userId, gender, userIds)
	var userQueue []string
	for _, userMatch := range usersMatch {
		userQueue = append(userQueue, strconv.Itoa(userMatch.Id))
	}
	return userQueue
}
