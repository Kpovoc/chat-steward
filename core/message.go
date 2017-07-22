package core

import "github.com/Kpovoc/JBot-Go/core/user"

type Message struct{
	Sender *user.User
	Content string
	//OriginatingLocation string // prvt Msg, general, coder, etc...
	//Recipient *user.User
}