package chatbot_user_handler

import (
	"fmt"
	"github.com/pa-pe/luca-control/src/storage"
	"github.com/pa-pe/luca-control/src/storage/model"
	"log"
	"strconv"
)

var cbServerErr = "oops, chatbot server error"

var functions = map[string]func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string){
	"getKeyboardOfOpenedShift": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return getKeyboardOfFlow(telegramStorage, 4)
	},
	"getKeyboardOfClosedShift": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return getKeyboardOfFlow(telegramStorage, 5)
	},
	"getLocationsKeyboard":     getLocationsKeyboard,
	"handleUserChooseLocation": handleUserChooseLocation,
	"handleRemainderProduct(OpenShift,FrameA)": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return handleRemainderProduct(1, telegramStorage, tgUser, msg, "open")
	},
	"handleRemainderProduct(OpenShift,FrameB)": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return handleRemainderProduct(2, telegramStorage, tgUser, msg, "open")
	},
	"handleRemainderProduct(OpenShift,Paper)": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return handleRemainderProduct(3, telegramStorage, tgUser, msg, "open")
	},
	"handleRemainderProduct(CloseShift,FrameA)": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return handleRemainderProduct(1, telegramStorage, tgUser, msg, "close")
	},
	"handleRemainderProduct(CloseShift,FrameB)": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return handleRemainderProduct(2, telegramStorage, tgUser, msg, "close")
	},
	"handleRemainderProduct(CloseShift,Paper)": func(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
		return handleRemainderProduct(3, telegramStorage, tgUser, msg, "close")
	}}

func Handle(telegramStorage storage.ITelegram, functionName string, tgUser *model.TgUser, msg string) (string, string) {
	if function, exists := functions[functionName]; exists {
		return function(telegramStorage, tgUser, msg)
	} else {
		log.Printf("chatbot_user_handler: Function '%s' not found!", functionName)
		return cbServerErr, ""
	}
}

func HandleServerError() (string, string) {
	return cbServerErr, ""
}

func getKeyboardOfFlow(telegramStorage storage.ITelegram, tgCbFlowId int) (string, string) {
	cbFlowSteps, err := telegramStorage.GetCbFlowAllSteps(tgCbFlowId)
	if err != nil {
		log.Printf("chatbot_user_handler: Error getting cbFlowSteps: %v", err)
		return cbServerErr, ""
	}

	return cbFlowSteps[0].Msg, cbFlowSteps[0].Keyboard
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

func handleRemainderProduct(srvsGoodsId int, telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string, shiftAction string) (string, string) {
	leftoverCount, err := strconv.Atoi(msg)
	if err != nil {
		return "Please enter just digit", ""
	}

	if strconv.Itoa(leftoverCount) != msg {
		return "Please enter just digit", ""
	}

	// get srvsShiftId from srvsEmployee
	srvsEmployeesList, err := telegramStorage.GetSrvsEmployeesList(fmt.Sprintf("id = %d", tgUser.SrvsEmployeesId))
	if err != nil {
		log.Printf("chatbot_user_handler: GetSrvsEmployeesList failed: %v", err)
		return cbServerErr, ""
	}
	srvsShiftId := srvsEmployeesList[0].SrvsShiftId

	// get SrvsLocationId from SrvsShifts
	srvsShifts, err := telegramStorage.GetSrvsShifts(fmt.Sprintf("id = %d", srvsShiftId))
	if err != nil {
		log.Printf("chatbot_user_handler: GetSrvsShift failed: %v", err)
		return cbServerErr, ""
	}
	srvsLocationId := srvsShifts[0].SrvsLocationId

	var srvsLeftover = model.SrvsLeftovers{
		SrvsShiftId:     srvsShiftId,
		SrvsLocationId:  srvsLocationId,
		SrvsGoodsId:     srvsGoodsId,
		SrvsEmployeesId: tgUser.SrvsEmployeesId,
	}

	if shiftAction == "open" {
		srvsLeftover.QuantityStart = leftoverCount
	} else if shiftAction == "close" {
		srvsLeftover.QuantityEnd = leftoverCount
	} else {
		log.Printf(`chatbot_user_handler: handleRemainderProduct unknown shiftAction="%s" expected "open" or "close"`, shiftAction)
		return cbServerErr, ""
	}

	_, err = telegramStorage.InsertSrvsLeftover(&srvsLeftover)
	if err != nil {
		return cbServerErr, ""
	}

	return "", ""
}

//func handleRemainderProductFrameA(telegramStorage storage.ITelegram, tgUser *model.TgUser, msg string) (string, string) {
//	return handleRemainderProduct(1, telegramStorage, tgUser, msg)
//}
