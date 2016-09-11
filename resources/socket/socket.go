package socket

import (
	"fmt"
	http "net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// TODO(Pholey): Origin checking
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Websocket failed to upgrade with error %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}

// Entry is The handler of the upgrade request
func Entry(c *gin.Context) {

	// We have determined they are a-OK to open a connection..
	// pass it to the handler
	wshandler(c.Writer, c.Request)
}
