package routes

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	def "../definitions"
)

var chatrooms = make([]*def.ChatRoom, 0)


func isChatActive(userId int) (*def.ChatRoom, bool) {
	for _, each := range chatrooms {
		if each.Users[0] == userId || each.Users[1] == userId {
			return each, true
		}
	}
	return nil, false
}

func ChatHandler(res http.ResponseWriter, req *http.Request) {
	currentConn, err := websocket.Upgrade(res, req, res.Header(),
		1024, 1024)
	if err != nil {
		http.Error(res, "could not open websocket connection",
		http.StatusBadRequest)
	}

	var users [2]int
	userId := req.URL.Query().Get("uid")
	chatWith := req.URL.Query().Get("cwid")

	chatConn, ok := isChatActive(int(chatWith))
	if !ok {
		users = append(users, userId)
		chatConn = def.ChatRoom{Conn: currentConn, Users: users}
		chatrooms = append(chatrooms, &chatConn)
	}

	go handleIO(&chatConn)
}

func handleIO(conn *def.ChatRoom) {
	broadcastMsg(conn, "Connected", "")

	for {
		payload := def.ChatMessage{}
		err := conn.ReadJSON(&payload)
		if err != nil {
			// remove connection
			break
		}
		broadcastMsg(conn, "MESSAGE", payload.Message)
	}
}

func broadcaseMsg(conn *def.ChatRoom, kind, message string) {
	conn.WriteJSON(def.ChatMessage {
		From: conn.Users[0]
		To: conn.Users[1]
		Type: kind,
		Message: message
	})
}







