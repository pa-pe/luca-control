package chatbot_user_handler

import (
	"fmt"
	"github.com/pa-pe/luca-control/src/storage"
	"github.com/pa-pe/luca-control/src/storage/model"
	"log"
)

var cbServerErr = "oops, chatbot server error"

var functions = map[string]func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string){
	"getLocationsKeyboard":     getLocationsKeyboard,
	"handleUserChooseLocation": handleUserChooseLocation,
}

func Handle(telegramStorage storage.ITelegram, functionName string, tgUser *model.TgUser, msg string) (string, string) {
	if function, exists := functions[functionName]; exists {
		return function(telegramStorage, tgUser, msg)
	} else {
		log.Printf("chatbot_user_handler: Function '%s' not found!", functionName)
	}
	return "", ""
}

func HandleServerError() (string, string) {
	return cbServerErr, ""
}

func getLocationsKeyboard(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
	_ = tgUser

	srvsLocationList, err := telegramStorage.GetSrvsLocationList("")
	if err != nil {
		log.Printf("chatbot_user_handler: GetSrvsLocationList failed: %v", err)
		return cbServerErr, ""
	}

	keyboard := ""
	for _, record := range srvsLocationList {
		if keyboard != "" {
			keyboard += "|"
		}
		keyboard += record.Name
	}

	return msg, keyboard
}

func handleUserChooseLocation(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
	_ = tgUser

	srvsLocationList, err := telegramStorage.GetSrvsLocationList(fmt.Sprintf(`name = "%s"`, msg))
	if err != nil {
		log.Printf("chatbot_user_handler: GetSrvsLocationList failed: %v", err)
		return cbServerErr, ""
	}

	if len(srvsLocationList) == 0 {
		return "Please tap locations button from menu", "func:getLocationsKeyboard"
	}

	// Start shift
	srvsShift := model.SrvsShifts{
		SrvsLocationId:  srvsLocationList[0].ID,
		SrvsEmployeesId: tgUser.SrvsEmployeesId,
	}
	_, err = telegramStorage.InsertSrvsShift(&srvsShift)
	if err != nil {
		log.Printf("chatbot_user_handler: InsertSrvsShift failed: %v", err)
		return cbServerErr, ""
	}

	err = telegramStorage.UpdateEmployeeSrvsShiftId(tgUser.SrvsEmployeesId, srvsShift.ID)
	if err != nil {
		log.Printf("chatbot_user_handler: UpdateSrvsShift failed: %v", err)
		return cbServerErr, ""
	}

	// return empty if handler pass userdata
	return "", ""
}
