package message

import (
	"time"

	"gitlab.com/Kpovoc/chat-steward/src/core/user"
)

type Message struct{
	Sender *user.User
	CreatedOn time.Time
	Content string
	//OriginatingLocation string // prvt Msg, general, coder, etc...
	//Recipient *user.User
}