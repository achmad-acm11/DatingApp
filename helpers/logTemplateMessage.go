package helpers

import (
	constants "DatingApp/constansts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

var timeElapsed *logTime

type StandartLog struct {
	TypeLayer constants.LayerName
	Name      constants.NameEntity
	NameFunc  string
}

func NewStandardLog(nameEntity constants.NameEntity, typeLayer constants.LayerName) *StandartLog {
	return &StandartLog{
		TypeLayer: typeLayer,
		Name:      nameEntity,
	}
}

func (log *StandartLog) StartFunction(request interface{}) {
	if gin.Mode() != gin.TestMode {
		timeElapsed = NewLogTime(time.Now())
		constants.Logger.WithField(string(constants.Request), request).Info(fmt.Sprintf("Start %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}

func (log *StandartLog) WarningFunction(message interface{}) {
	if gin.Mode() != gin.TestMode {
		constants.Logger.WithField("message", message).Info(fmt.Sprintf("Warning %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}

func (log *StandartLog) EndFunction(response interface{}) {
	if gin.Mode() != gin.TestMode {
		constants.Logger.WithFields(logrus.Fields{
			string(constants.Response):    response,
			string(constants.TimeElapsed): timeElapsed.GetTimeSince(),
		}).Info(fmt.Sprintf("End %s %s %s", log.TypeLayer, log.Name, log.NameFunc))
	}
}
