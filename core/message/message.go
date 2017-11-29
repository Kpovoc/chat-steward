package message

import (
	"github.com/Kpovoc/JBot-Go/core/user"
	"time"
)

type Message struct{
	Sender *user.User
	CreatedOn time.Time
	Content string
	//OriginatingLocation string // prvt Msg, general, coder, etc...
	//Recipient *user.User
}