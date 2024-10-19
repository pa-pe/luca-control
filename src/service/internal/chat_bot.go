package internal

import (
	"github.com/pa-pe/luca-control/src/storage"
	"github.com/pa-pe/luca-control/src/storage/model"
	"log"
)

type ChatBotImpl struct {
	telegramStorage     storage.ITelegram
	chatBotUserInserted bool
}

// Handle returns Msg string and func which needs to be called in case of successful sending the Msg
func (c *ChatBotImpl) Handle(botTgUser model.TgUser, tgUser model.TgUser, tgMsg model.TgMsg) (string, func(tgId int64)) {

	_, err := c.telegramStorage.InsertMsg(&tgMsg)
	if err != nil {
		log.Print("ChatBot Handle problem")
		return "", nil
	}

	// check once for existence chatBotUser in db
	if c.chatBotUserInserted == false {
		err = c.telegramStorage.CreateUserIfNotExist(&botTgUser)
		if err != nil {
			log.Print("ChatBot Handle problem with db insert chatBotTgUser")
			return "", nil
		}
		c.chatBotUserInserted = true
	}

	err = c.telegramStorage.CreateUserIfNotExist(&tgUser)
	if err != nil {
		log.Print("ChatBot Handle problem with db insert tgUser")
		return "", nil
	}

	answerMsg := c.echo(tgMsg.Text)

	// finish if no answer msg
	if answerMsg == "" {
		return "", nil
	}

	tgMsgOut := model.TgMsg{
		ChatID:     tgMsg.ChatID,
		TgUserID:   botTgUser.ID,
		Text:       answerMsg,
		IsOutgoing: 1,
		//		AddedTimestamp: time.Now().Unix(),
	}

	_, err = c.telegramStorage.InsertMsg(&tgMsgOut)
	if err != nil {
		log.Print("ChatBot Handle problem with insert tgMsgOut")
		return "", nil
	}

	executeAfterSent := func(tgId int64) {
		//		log.Print("sent: " + answerMsg)
		tgMsgOut.TgID = tgId
		_ = c.telegramStorage.UpdateTgOutMsgIdAfterSend(&tgMsgOut)
	}

	return answerMsg, executeAfterSent
}

func (c *ChatBotImpl) echo(msg string) string {
	return "Echo: " + msg
}

func NewChatBotService(telegramStorage storage.ITelegram) *ChatBotImpl {
	return &ChatBotImpl{
		telegramStorage:     telegramStorage,
		chatBotUserInserted: false,
	}
}
