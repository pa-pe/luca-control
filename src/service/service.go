package service

import (
	"github.com/pa-pe/luca-control/src/service/internal"
	"github.com/pa-pe/luca-control/src/storage"
	"github.com/pa-pe/luca-control/src/storage/model"
)

type Services struct {
	ChatBotMsgRouter IChatBotMsg
}

type IChatBotMsg interface {
	Handle(botUser model.TgUser, tgUser model.TgUser, tgMsg model.TgMsg) (string, func(tgId int64))
}

func NewServices(storage *storage.Storages) *Services {
	return &Services{
		ChatBotMsgRouter: internal.NewChatBotService(storage.Telegram),
	}
}
