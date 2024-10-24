package internal

import (
	"fmt"
	"github.com/pa-pe/luca-control/src"
	"github.com/pa-pe/luca-control/src/service/internal/chatbot_user_handler"
	"github.com/pa-pe/luca-control/src/storage"
	"github.com/pa-pe/luca-control/src/storage/model"
	"log"
	"strings"
	"unicode"
)

type ChatBotImpl struct {
	TelegramStorage     storage.ITelegram
	tgBot               *src.BotImpl
	chatBotUserInserted bool
}

// Handle returns Msg string and func which needs to be called in case of successful sending the Msg
func (c *ChatBotImpl) Handle(botTgUser model.TgUser, tgUser model.TgUser, tgMsg model.TgMsg) (string, string, func(tgId int64)) {

	_, err := c.TelegramStorage.InsertMsg(&tgMsg)
	if err != nil {
		log.Print("ChatBot Handle problem")
		return "", "", nil
	}

	// check once for existence chatBotUser in db
	if c.chatBotUserInserted == false {
		_, err = c.TelegramStorage.CreateUserIfNotExist(&botTgUser)
		if err != nil {
			log.Print("ChatBot Handle problem with db insert chatBotTgUser")
			return "", "", nil
		}
		c.chatBotUserInserted = true
	}

	isUserCreated, err := c.TelegramStorage.CreateUserIfNotExist(&tgUser)
	if err != nil {
		log.Print("ChatBot Handle problem with db insert tgUser")
		return "", "", nil
	}
	if isUserCreated {
		c.onNewTgUser(&tgUser)
	}

	//	answerMsg := c.echo(tgMsg.Text)
	answerMsg, keyboardStr := c.msgRouter(&tgMsg)
	answerMsg, keyboardStr = c.UnFuncKeyboard(&tgUser, answerMsg, keyboardStr)

	// finish if no answer msg
	if answerMsg == "" {
		if keyboardStr != "" {
			log.Printf("Trying to send keyboard '%s' without msg", keyboardStr)
		}
		return "", "", nil
	}

	tgMsgOut := model.TgMsg{
		ChatID:     tgMsg.ChatID,
		TgUserID:   botTgUser.ID,
		Text:       answerMsg,
		IsOutgoing: 1,
		//		AddedTimestamp: time.Now().Unix(),
	}

	// insert keyboard info for db store
	if keyboardStr != "" {
		tgMsgOut.Text = answerMsg + "\n\nkeyboard:\n" + keyboardStr
	}

	_, err = c.TelegramStorage.InsertMsg(&tgMsgOut)
	if err != nil {
		log.Print("ChatBot Handle problem with insert tgMsgOut")
		return "", "", nil
	}

	// restore msg without keyboard info
	tgMsgOut.Text = answerMsg

	executeAfterSent := func(tgId int64) {
		//		log.Print("sent: " + answerMsg)
		tgMsgOut.TgID = tgId
		_ = c.TelegramStorage.UpdateTgOutMsgIdAfterSend(&tgMsgOut)
	}

	return answerMsg, keyboardStr, executeAfterSent
}

func (c *ChatBotImpl) msgRouter(tgMsg *model.TgMsg) (string, string) {
	msg := tgMsg.Text
	msgLc := strings.ToLower(msg)

	var builder strings.Builder
	for _, r := range msgLc {
		if unicode.IsLetter(r) {
			builder.WriteRune(r)
		}
	}
	msgOnlyLetters := builder.String()
	_ = msgOnlyLetters

	tgUser, err := c.TelegramStorage.FindUserById(tgMsg.TgUserID)
	if err != nil {
		log.Fatalf("tgUser not found! WHYYYY???? error: %v", err)
	}

	if tgUser.ChatbotPermit != 0 {
		if tgUser.TgCbFlowStepId != 0 {
			return c.handleUserContinueFlow(tgUser, msg)
		}

		if msg == "Start shift" {
			return c.handleUserStartFlow(tgUser, 1)
		} else {
			return c.handleUserStartFlow(tgUser, 4)
		}
	}

	if msgLc == "/start" {
		//answer := HandleCmdStart()
		answer := "Hello! Please wait for permission to continue."
		return answer, ""
		//} else if msgOnlyLetters == "hi" {
		//	return "Hello!", ""
		//} else if msgOnlyLetters == "hello" {
		//	return "Hi!", ""
		//} else if msgLc == "kb" {
		//	return "try!", "Кнопка 1|Кнопка 2|\nКнопка 3#hide"
		//} else if msgLc == "kb2" {
		//	return "try!", "Remove KB|Hide KB"
		//} else if msg == "Remove KB" {
		//	return "removed", "remove"
		//} else if msg == "Hide KB" {
		//	return "...", ""
		//} else if msg == "users" {
		//	answer := ""
		//tgUsers, _ := c.TelegramStorage.FindUsersByCustomQuery("")
		//for _, record := range *tgUsers {
		//	answer = answer + record.UserName + "\n"
		//}
		//return answer, ""
		//} else if msg == "msg" {
		//	c.tgBot.TgController.SendMessage(tgMsg.ChatID, ";)", "")
		//	return "done", ""
	}

	answer := "You do not have permission to work with the chatbot.\nPlease wait for permission to continue."
	return answer, "remove"

	//answers := make([]string, 0)
	//answers = append(answers, "0_o", "o_0", "o_o", "0_0")
	////return "0_o", ""
	//return answers[rand.Intn(len(answers))], ""
}

func (c *ChatBotImpl) handleUserStartFlow(tgUser *model.TgUser, tgCbFlowId int) (string, string) {
	tgCbFlowAllSteps, err := c.TelegramStorage.GetCbFlowAllSteps(tgCbFlowId)
	if err != nil {
		log.Printf("ChatBot handleUserStartFlow: GetCbFlowAllSteps error: %v", err)
		return chatbot_user_handler.HandleServerError()
	}

	// save current step just if exist next step or handler
	if len(tgCbFlowAllSteps) > 1 || tgCbFlowAllSteps[0].HandlerName != "" {
		err = c.TelegramStorage.UpdateTgUserFlowStep(tgUser.ID, tgCbFlowAllSteps[0].ID)
		if err != nil {
			log.Printf("ChatBot handleUserStartFlow: UpdateTgUserFlowStep error: %v", err)
			return chatbot_user_handler.HandleServerError()
		}
	}

	return tgCbFlowAllSteps[0].Msg, tgCbFlowAllSteps[0].Keyboard
}

func (c *ChatBotImpl) handleUserContinueFlow(tgUser *model.TgUser, msg string) (string, string) {
	currentTgCbFlowStep, err := c.TelegramStorage.GetCbFlowStep(tgUser.TgCbFlowStepId)
	if err != nil {
		log.Printf("ChatBot handleUserContinueFlow: GetCbFlowStep error: %v", err)
		return chatbot_user_handler.HandleServerError()
	}

	if currentTgCbFlowStep.HandlerName != "" {
		msg, keyboard := chatbot_user_handler.Handle(c.TelegramStorage, currentTgCbFlowStep.HandlerName, tgUser, msg)
		if msg != "" || keyboard != "" {
			// chatbot_user_handler return msg if it can't recognize user msg or other error
			return msg, keyboard
		}
	}

	// continue next steps
	nextTgCbFlowStep, err := c.TelegramStorage.GetNextCbFlowStep(tgUser.TgCbFlowStepId)
	if err != nil {
		log.Printf("ChatBot handleUserContinueFlow: GetNextCbFlowStep error: %v", err)
		return chatbot_user_handler.HandleServerError()
	}

	if nextTgCbFlowStep == nil {
		// flow finished
		err := c.TelegramStorage.UpdateTgUserFlowStep(tgUser.ID, 0)
		if err != nil {
			log.Printf("ChatBot handleUserContinueFlow: UpdateTgUserFlowStep error: %v", err)
			return chatbot_user_handler.HandleServerError()
		}
		return "", "" // keep silent in case of completion
	}

	return nextTgCbFlowStep.Msg, nextTgCbFlowStep.Keyboard
}

//func (c *ChatBotImpl) handleUserFlowNextStep(tgUser *model.TgUser) (string, string) {
//	return "", ""
//}

func (c *ChatBotImpl) echo(msg string) string {
	return "Echo: " + msg
}

//func HandleCmdStart() string {
//	answer := "Hello! Please wait for permission to continue."
//	return answer
//}

func (c *ChatBotImpl) onNewTgUser(newTgUser *model.TgUser) {
	foundTgUsers, _ := c.TelegramStorage.FindUsersByCustomQuery("id = 568876500")
	for _, foundTgUser := range *foundTgUsers {
		msg := "Admin: New user connected to telegram bot, UserName=" + newTgUser.UserName + ", FirstName=" + newTgUser.FirstName
		c.tgBot.TgController.SendMessage(foundTgUser.ID, msg, "")
	}
}

func (c *ChatBotImpl) getLocationsKeyboard() string {
	var keyboard string

	return keyboard
}

func (c *ChatBotImpl) UnFuncKeyboard(tgUser *model.TgUser, msg string, keyboard string) (string, string) {
	if len(keyboard) > 5 && keyboard[0:5] == "func:" {
		fmt.Printf("UnFuncKeyboard: '%s' -> '%s'\n ", keyboard, keyboard[5:])
		function := keyboard[5:]
		msg, keyboard = chatbot_user_handler.Handle(c.TelegramStorage, function, tgUser, msg)
	}

	return msg, keyboard
}

func NewChatBotService(telegramStorage storage.ITelegram, tgBot *src.BotImpl) *ChatBotImpl {
	return &ChatBotImpl{
		TelegramStorage:     telegramStorage,
		tgBot:               tgBot,
		chatBotUserInserted: false,
	}
}
