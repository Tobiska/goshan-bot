package models

type IncomingMessage struct {
	ChatID          int64
	UserID          int64
	Username        string
	UsernameDisplay string
	Text            string
	IsCallback      bool
	CallbackMsgID   string
}
