package storage

import (
	"github.com/pa-pe/luca-control/src/storage/internal"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
)

type Storages struct {
	Telegram ITelegram
}

type ITelegram interface {
	FindUserById(userID int64) (*model.TgUser, error)
	FindUsersByCustomQuery(where string) (*[]model.TgUser, error)
	CreateUserIfNotExist(tgUser *model.TgUser) (bool, error)
	InsertMsg(tgMsg *model.TgMsg) (int64, error)
	UpdateTgOutMsgIdAfterSend(tgMsg *model.TgMsg) error
	GetCbFlow(tgCbFlowId int) (*[]model.TgCbFlowStep, error)
}

func NewStorages(db *gorm.DB) *Storages {
	return &Storages{
		Telegram: internal.NewTelegramStorage(db),
	}
}
