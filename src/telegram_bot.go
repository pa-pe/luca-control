package src

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pa-pe/luca-control/src/service"
	"github.com/pa-pe/luca-control/src/storage/model"
	"log"
	"strings"
)

type BotImpl struct {
	token    string
	services *service.Services
}

func NewTelegramBot(token string, services *service.Services) *BotImpl {
	return &BotImpl{token: token, services: services}
}

func (bot *BotImpl) ListenAndServ() {
	tgController, err := NewTelegramController(bot)
	if err != nil {
		//		log.Fatalf("Error creating Telegram tgController: %v", err)
		log.Fatalf("Error creating Telegram tgController: %v", err)

	}

	//	go func() {
	if err := tgController.checkConnection(); err != nil {
		//			log.Fatalf("Failed to connect to Telegram: %v", err)
		log.Fatalf("Failed to connect to Telegram: %v", err)
	}
	log.Print("Telegram bot started with UserName=" + tgController.BotInfo.UserName)

	if err := tgController.ListenForMessages(); err != nil {
		//			log.Fatalf("Error listening for messages: %v", err)
		log.Fatalf("Error listening for messages: %v", err)
	}
	//	}()
}

type TelegramController struct {
	bot     *BotImpl
	botAPI  *tgbotapi.BotAPI
	BotInfo *tgbotapi.User
}

func NewTelegramController(bot *BotImpl) (*TelegramController, error) {
	newBotAPI, err := tgbotapi.NewBotAPI(bot.token)
	if err != nil {
		return nil, err
	}
	newBotAPI.Debug = true
	//	newBotAPI.Self.

	return &TelegramController{
		bot:    bot,
		botAPI: newBotAPI,
	}, nil
}

func (c *TelegramController) checkConnection() error {
	return c.fetchBotInfo()
}

func (c *TelegramController) fetchBotInfo() error {
	botUser, err := c.botAPI.GetMe()
	if err != nil {
		return err
	}

	c.BotInfo = &botUser
	//	log.Printf("Bot Info: TgID=%d, Username=%s, FirstName=%s, LastName=%s, LanguageCode=%s", c.BotInfo.TgID, c.BotInfo.UserName, c.BotInfo.FirstName, c.BotInfo.LastName, c.BotInfo.LanguageCode)
	return nil
}

func (c *TelegramController) ListenForMessages() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.botAPI.GetUpdatesChan(u)

	botUser := model.TgUser{
		ID:        c.BotInfo.ID,
		UserName:  c.BotInfo.UserName,
		FirstName: c.BotInfo.FirstName,
		LastName:  c.BotInfo.LastName,
	}

	for update := range updates {
		if update.Message != nil {
			user := model.TgUser{
				ID:           update.Message.From.ID,
				UserName:     update.Message.From.UserName,
				FirstName:    update.Message.From.FirstName,
				LastName:     update.Message.From.LastName,
				LanguageCode: update.Message.From.LanguageCode,
			}

			msg := model.TgMsg{
				TgID:     int64(update.Message.MessageID),
				Text:     update.Message.Text,
				ChatID:   update.Message.Chat.ID,
				TgUserID: update.Message.From.ID,
				Date:     int64(update.Message.Date),
				ReplyToMessageID: func() int64 {
					if update.Message.ReplyToMessage != nil {
						return int64(update.Message.ReplyToMessage.MessageID)
					}
					return 0
				}(),
				//				AddedTimestamp: time.Now().Unix(),
			}

			// if received the sticker
			if update.Message.Sticker != nil {
				// convert to text emoji
				if update.Message.Sticker.FileUniqueID != "" {
					msg.Text = "sticker converted to emoji: " + update.Message.Sticker.Emoji
				}
			}

			fmt.Printf("\ntg: msgId=%d from=%s[%d] text=%s\n", update.Message.MessageID, update.Message.From.UserName, update.Message.From.ID, msg.Text)

			//			responseMsg := handler.ChatBotMsgProcess(msg.Text)
			responseMsg, keyboardStr, executeAfterSent := c.bot.services.ChatBotMsgRouter.Handle(botUser, user, msg)

			// sending a response if the handler gave one
			if responseMsg != "" {

				fmt.Printf("tg: answer text=%s\n\n", responseMsg)
				response := tgbotapi.NewMessage(msg.ChatID, responseMsg)

				if keyboardStr == "remove" {
					removeKeyboard := tgbotapi.NewRemoveKeyboard(true) // true означает, что клавиатура будет удалена для всех пользователей
					response.ReplyMarkup = removeKeyboard
				} else if keyboardStr != "" {
					keyboard := loadKeyboardFromTxt(keyboardStr)
					response.ReplyMarkup = keyboard
				}

				sent, err := c.botAPI.Send(response)
				if err != nil {
					return err
				}

				if executeAfterSent != nil {
					executeAfterSent(int64(sent.MessageID))
				}
			}
		}
	}

	return nil
}

func loadKeyboardFromTxt(input string) tgbotapi.ReplyKeyboardMarkup {
	// Определение, нужно ли включить one_time_keyboard
	oneTimeKeyboard := strings.Contains(input, "#hide")

	// Удаляем маркер #hide из строки, если он присутствует
	input = strings.ReplaceAll(input, "#hide", "")

	// Разделяем строки по символу \n
	rows := strings.Split(input, "\n")

	// Формируем клавиатуру
	var keyboardRows [][]tgbotapi.KeyboardButton
	for _, row := range rows {
		// Разделяем кнопки по символу |
		buttons := strings.Split(row, "|")

		// Создаем строку с кнопками
		var keyboardButtons []tgbotapi.KeyboardButton
		for _, text := range buttons {
			text = strings.TrimSpace(text)
			if text != "" {
				keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(text))
			}
		}

		// Добавляем строку кнопок, если она не пустая
		if len(keyboardButtons) > 0 {
			keyboardRows = append(keyboardRows, keyboardButtons)
		}
	}

	// Создание ReplyKeyboardMarkup
	keyboard := tgbotapi.NewReplyKeyboard(keyboardRows...)
	keyboard.OneTimeKeyboard = oneTimeKeyboard
	keyboard.ResizeKeyboard = true

	return keyboard
}
