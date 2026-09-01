package handlers

import (
	"net/http"

	appws "infrapilot/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *appws.Hub) gin.HandlerFunc {

	return func(c *gin.Context) {

		conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			return
		}

		client := &appws.Client{
			Conn: conn,
			Send: make(chan []byte, 256), // Ensure non-blocking buffered channel
		}

		hub.Register <- client

		go func() {
			for msg := range client.Send {
				client.Conn.WriteMessage(websocket.TextMessage, msg)
			}
		}()
	}
}
