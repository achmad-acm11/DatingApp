package tests

import (
	"DatingApp/entities"
	"DatingApp/helpers"
	"database/sql/driver"
)

var MatchQueueCols = []string{"id", "user_id", "pass_count", "like_count", "current_state", "user_queue", "date"}

var MatchQueueListDummy = []entities.MatchQueue{
	entities.MatchQueue{
		UserId:       1,
		PassCount:    0,
		LikeCount:    0,
		CurrentState: 1,
		UserQueue:    "2|3|4",
		Date:         helpers.GetLocalDateNow(),
	}}

func MappingMatchQueueStore(item entities.MatchQueue, numId int) []driver.Value {
	var values []driver.Value

	values = append(values, numId)
	values = append(values, item.UserId)
	values = append(values, item.PassCount)
	values = append(values, item.LikeCount)
	values = append(values, item.CurrentState)
	values = append(values, item.UserQueue)
	values = append(values, item.Date)

	return values
}
