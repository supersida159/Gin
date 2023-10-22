package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Upgrade this connection to WebSocket
func Upgrade(c *gin.Context) (*websocket.Conn, error) {
	// Allow WebSocket connections from any origin
	// upgrader.CheckOrigin = func(r *http.Request) bool {
	// 	// List of allowed origins
	// 	allowedOrigins := []string{"http://example.com", "https://example.org"}

	// 	// Get the origin header from the request
	// 	origin := r.Header.Get("Origin")

	// 	// Check if the origin is in the list of allowed origins
	// 	for _, allowedOrigin := range allowedOrigins {
	// 		if origin == allowedOrigin {
	// 			return true
	// 		}
	// 	}

	// 	return false
	// }
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return conn, nil
}
