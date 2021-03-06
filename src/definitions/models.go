package definitions

import (
	"github.com/gorilla/websocket"
)

type Request struct {
	Meta map[string]interface{} `json:"meta"`
	Data []map[string]interface{} `json:"data"`
}

type Response struct {
	Status int8 `json:"status"`
	Message string `json:"message"`
	Data []map[string]interface{} `json:"data"`
}

type CU struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedOn string `json:"createdOn"`
	Status string `json:"status"`
}


type Meta struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	ReqType string `json:"reqType"`
	Extra1 string `json:"extra1"`
	Extra2 string `json:"extra2"`
	Extra3 string `json:"extra3"`
	Extra4 string `json:"extra4"`
}

type ChatRoom struct {
	*websocket.Conn
	Users [2]int
}

type ChatMessage struct {
	From string
	To string
	Type string
	Message string
}
