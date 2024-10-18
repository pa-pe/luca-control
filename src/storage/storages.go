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
	//	FindFindUserById(ctx context.Context, userID int64) (*model.TgUser, error)
	FindUserById(userID int64) (*model.TgUser, error)
	CreateUserIfNotExist(tgUser *model.TgUser) error
	InsertMsg(tgMsg *model.TgMsg) (int64, error)
	UpdateTgOutMsgIdAfterSend(tgMsg *model.TgMsg) error
}

func NewStorages(db *gorm.DB) *Storages {
	return &Storages{
		Telegram: internal.NewTelegramStorage(db),
	}
}
