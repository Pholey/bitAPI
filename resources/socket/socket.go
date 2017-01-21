package socket

import (
	"fmt"
	http "net/http"

	// "github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
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

// func middleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
//
// 	}
// }

// Entry is The handler of the upgrade request
var Entry = echo.WrapHandler(http.HandlerFunc(wshandler))
