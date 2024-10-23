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
			log.Printf("User search error: %v", err)
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
			log.Printf("User search error: %v", err)
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
	db := c.DB
	//	if err := db.Create(&tgMsg).Error; err != nil {
	if result := db.Create(&tgMsg); result.Error != nil {
		log.Printf("Error inserting msg: %v", result.Error)
		return tgMsg.InternalID, result.Error
	}
	return 0, nil
}

func (c *TelegramImpl) UpdateTgOutMsgIdAfterSend(tgMsgOut *model.TgMsg) error {
	//	log.Printf("try to update tgMsgOut.internal_id=%d, set tg_id=%d, txt=%s\n", tgMsgOut.InternalID, tgMsgOut.TgID, tgMsgOut.Text)
	db := c.DB
	if err := db.Model(&tgMsgOut).Where("internal_id", tgMsgOut.InternalID).Update("tg_id", tgMsgOut.TgID).Update("date", time.Now().Unix()).Error; err != nil {
		return err
	}
	return nil
}

func (c *TelegramImpl) GetCbFlow(tgCbFlowId int) (*[]model.TgCbFlowStep, error) {
	var tgCbFlow *[]model.TgCbFlowStep
	if err := c.DB.Where("tgCbFlowId = ?", tgCbFlowId).Find(&tgCbFlow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//		fmt.Printf("UserID=%d not found\n", userID)
			return nil, nil
		} else {
			log.Printf("User search error: %v", err)
			return nil, err
		}
	}

	return tgCbFlow, nil
}

func NewTelegramStorage(db *gorm.DB) *TelegramImpl {
	return &TelegramImpl{
		DB: db,
	}
}
