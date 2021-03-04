package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type ChatRoom struct {
	*websocket.Conn
	Username string
}

type ChatMessage struct {
	From string
	To string
	Type string
	Message string
}

var chatrooms = make([]*ChatRoom, 0)

func main() {
	http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		currentConn, err := websocket.Upgrade(w, r, w.Header(),
			1024, 1024)
		if err != nil {
			http.Error(w, "Could not open websocket connection",
				http.StatusBadRequest)
		}

		username := r.URL.Query().Get("username")
		conn := ChatRoom{Conn: currentConn, Username: username}

		chatrooms = append(chatrooms, &conn)

		go handleIO(&conn, chatrooms)
	})
	fmt.Println("Server starting at :8080")
	http.ListenAndServe(":8080", nil)
}

func handleIO(conn *ChatRoom, chatrooms []*ChatRoom) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	broadcastMessage(conn, "New USER", "")

	for {
		payload := ChatMessage{}
		err := conn.ReadJSON(&payload)
		if err != nil {
			// remove connection 
			continue
		}
		broadcastMessage(conn, "MESSAGE", payload.Message)
	}
}

func broadcastMessage(conn *ChatRoom, kind, message string) {
	for _, each := range chatrooms{
		if each == conn {
			continue
		}

		each.WriteJSON(ChatMessage {
			From: conn.Username,
			To: "",
			Type: kind,
			Message: message,
		})
	}
}









