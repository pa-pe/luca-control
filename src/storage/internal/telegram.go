package internal

import (
	"errors"
	"github.com/pa-pe/luca-control/src/storage/model"
	"gorm.io/gorm"
	"log"
	"time"
)

type TelegramImpl struct {
	DB *gorm.DB
}

func (c *TelegramImpl) FindUserById(userID int64) (*model.TgUser, error) {
	//func (c *telegramImpl) FindUserById(ctx context.Context, userID int64) (*model.TgUser, error) {
	//	db := c.DB.WithContext(ctx)
	db := c.DB
	//	m := &model.TgUser{}
	//	db.Raw("Select 'Hello world' as title").Scan(m)

	var existingUser model.TgUser
	//	if err := db.Where("UserID = ?", userID).First(&existingUser).Error; err != nil {
	if err := db.Where("id = ?", userID).Take(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("FindUserById error: %v", err)
			return nil, err
		}
	}
	//fmt.Printf("UserID=%d found, UserName=%s\n", userID, existingUser.UserName)

	return &existingUser, nil
}

func (c *TelegramImpl) FindUsersByCustomQuery(where string) (*[]model.TgUser, error) {
	db := c.DB
	//	m := &model.TgUser{}

	var existingUsers *[]model.TgUser
	if err := db.Where(where).Find(&existingUsers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("FindUsersByCustomQuery error: %v", err)
			return nil, err
		}
	}

	return existingUsers, nil
}

func (c *TelegramImpl) CreateUserIfNotExist(tgUser *model.TgUser) (bool, error) {
	// check if record exists before adding
	existingTgUser, err := c.FindUserById(tgUser.ID)
	if err != nil {
		return false, err
	}

	db := c.DB

	if existingTgUser == nil {
		if err := db.Create(&tgUser).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// InsertMsg returns inserted ID and error
func (c *TelegramImpl) InsertMsg(tgMsg *model.TgMsg) (int64, error) {
	//	if err := db.Create(&tgMsg).Error; err != nil {
	if result := c.DB.Create(&tgMsg); result.Error != nil {
		log.Printf("InsertMsg error: %v", result.Error)
		return tgMsg.InternalID, result.Error
	}
	return 0, nil
}

func (c *TelegramImpl) UpdateTgOutMsgIdAfterSend(tgMsgOut *model.TgMsg) error {
	//	log.Printf("try to update tgMsgOut.internal_id=%d, set tg_id=%d, txt=%s\n", tgMsgOut.InternalID, tgMsgOut.TgID, tgMsgOut.Text)
	if err := c.DB.Model(&tgMsgOut).Where("internal_id", tgMsgOut.InternalID).Update("tg_id", tgMsgOut.TgID).Update("date", time.Now().Unix()).Error; err != nil {
		return err
	}
	return nil
}

func (c *TelegramImpl) GetCbFlowAllSteps(tgCbFlowId int) ([]model.TgCbFlowStep, error) {
	var tgCbFlowAllSteps []model.TgCbFlowStep
	if err := c.DB.Where("tg_cb_flow_id = ?", tgCbFlowId).Order("row_order").Find(&tgCbFlowAllSteps).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("GetCbFlowAllSteps error: %v", err)
			return nil, err
		}
	}

	return tgCbFlowAllSteps, nil
}

func (c *TelegramImpl) GetCbFlowStep(tgCbFlowStepId int) (*model.TgCbFlowStep, error) {
	var tgCbFlowStep *model.TgCbFlowStep
	if err := c.DB.Where("id = ?", tgCbFlowStepId).Find(&tgCbFlowStep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("GetCbFlowStep error: %v", err)
			return nil, err
		}
	}

	return tgCbFlowStep, nil
}

func (c *TelegramImpl) GetNextCbFlowStep(tgCbFlowStepId int) (*model.TgCbFlowStep, error) {
	currentTgCbFlowStep, err := c.GetCbFlowStep(tgCbFlowStepId)
	if err != nil {
		log.Printf("GetNextCbFlowStep error: %v", err)
		return nil, err
	}

	var nextTgCbFlowStep *model.TgCbFlowStep
	if err := c.DB.Where("tg_cb_flow_id = ?", currentTgCbFlowStep.TgCbFlowId).
		Where("row_order > ?", currentTgCbFlowStep.RowOrder).
		Order("row_order").
		First(&nextTgCbFlowStep).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("TgCbFlowStep search error: %v", err)
			return nil, err
		}
	}

	return nextTgCbFlowStep, nil
}

func (c *TelegramImpl) UpdateTgUserFlowStep(tgUserId int64, tgCbFlowStepId int) error {
	if err := c.DB.Model(model.TgUser{}).Where("id", tgUserId).Update("tg_cb_flow_step_id", tgCbFlowStepId).Error; err != nil {
		return err
	}
	return nil
}

func (c *TelegramImpl) GetSrvsLocationList(where string) ([]model.SrvsLocationList, error) {
	var srvsLocationList []model.SrvsLocationList

	if err := c.DB.Where(where).Find(&srvsLocationList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("GetSrvsLocationList search error: %v", err)
			return nil, err
		}
	}

	return srvsLocationList, nil
}
func (c *TelegramImpl) InsertSrvsShift(srvsShift *model.SrvsShifts) (int, error) {
	if result := c.DB.Create(&srvsShift); result.Error != nil {
		log.Printf("InsertSrvsShift error: %v", result.Error)
		return srvsShift.ID, result.Error
	}
	return 0, nil
}

func (c *TelegramImpl) UpdateEmployeeSrvsShiftId(srvsEmployeeId int, srvsShiftId int) error {
	if err := c.DB.Model(model.SrvsEmployeesList{}).Where("id", srvsEmployeeId).Update("srvs_shift_id", srvsShiftId).Error; err != nil {
		return err
	}
	return nil
}

func NewTelegramStorage(db *gorm.DB) *TelegramImpl {
	return &TelegramImpl{
		DB: db,
	}
}
