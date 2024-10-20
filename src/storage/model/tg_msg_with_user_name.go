package model

type TgMsgWithUserName struct {
	InternalID       int64  `json:"internal_id"`
	TgID             int64  `json:"tg_id"`
	TgUserID         int64  `json:"tg_user_id"`
	UserName         string `json:"user_name"`
	ChatID           int64  `json:"chat_id"`
	ReplyToMessageID int64  `json:"reply_to_message_id"`
	IsOutgoing       byte   `json:"is_outgoing"`
	Text             string `json:"text"`
	Date             int64  `json:"date"`
	AddedTimestamp   int64  `json:"added_timestamp"`
}
