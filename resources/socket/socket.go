package socket

import (
	"fmt"
	http "net/http"

	"golang.org/x/net/websocket"
	"github.com/labstack/echo"
	"github.com/speps/go-hashids"
	mw "github.com/Pholey/bitAPI/resources/middleware"
)

// var wsupgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
//
// 	// TODO(Pholey): Origin checking
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// func wshandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := wsupgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Websocket failed to upgrade with error %+v", err)
// 		return
// 	}
//
// 	for {
// 		t, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		conn.WriteMessage(t, msg)
// 	}
// }

// func middleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
//
// 	}
// }

var socketMaxId int64 = 0
var socketHd = &hashids.HashIDData{hashids.DefaultAlphabet, 64, "socket"}

var msgMaxId int64 = 0
var msgHd = &hashids.HashIDData{hashids.DefaultAlphabet, 64, "message"}

type Socket struct {
	id       string
	channels map[string]*Channel
	addCh    chan *Channel
	delCh    chan *Channel
	sendCh   chan *Message
	doneCh   chan bool
	errCh    chan error
}

func NewSocket() (*Socket, error) {
	channels := make(map[string]*Channel)
	addCh := make(chan *Channel)
	delCh := make(chan *Channel)
	sendCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	h := hashids.NewWithData(socketHd)
	socketMaxId++
  e, err := h.EncodeInt64([]int64{socketMaxId})
	if err != nil {
		return nil, err
	}

	return &Socket{e, channels, addCh, delCh, sendCh, doneCh, errCh}, nil
}

func (s *Socket) HasRecipient(id string) bool {
	if s.id == id {
		return true
	}
	_, ok := s.channels[id]
	return ok
}

func (s *Socket) Add(c *Channel) {
	// queue, ok := s.outbox[c.id]
	// if ok {
	// 	for _, msg := range(queue) {
	// 		err := c.Write(msg)
	// 		if (err != nil) {
	// 			break
	// 		}
	// 	}
	// }
	s.addCh <- c
}

func (s *Socket) Del(c *Channel) {
	s.delCh <- c
}

func (s *Socket) Done() {
	s.doneCh <- true
}

func (s *Socket) Err(err error) {
	s.errCh <- err
}

func (s *Socket) send(msg *Message) error {
	if s.id != msg.To {
		channel, ok := s.channels[msg.To]
		if ok {
			channel.Write(msg)
		}
		// else {
		// 	outBuf, ok := s.outbox[msg.To]
		// 	if !ok {
		// 		outBuf = make([]*Message)
		// 	}
		// 	s.outbox[msg.To] = append(outBuf, msg)
		// }
	} // Should have an else here for messages to the server itself

	return nil
}

func socketLoop(s *Socket) error {
	for {
		select {
		case c := <-s.addCh:
			s.channels[c.id] = c
		case c := <-s.delCh:
			delete(s.channels, c.id)
		case msg := <-s.sendCh:
			s.send(msg)
		case err := <-s.errCh:
			fmt.Println("Error:", err.Error())
			return err
		case <-s.doneCh:
			return nil
		}
	}
}

func (s *Socket) Listen(c echo.Context) error {
	id := c.Param("id")

	websocket.Handler(func(ws *websocket.Conn) {
		closeSocket := func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}

		cu := c.(*mw.ContextWithUser)

		channel, err := NewChannel(ws, s, cu.UserId, id)
		if err != nil {
			closeSocket()
			return
		}
		go socketLoop(s)
		s.Add(channel)
		channel.Listen()
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}

func (s *Socket) GetChannels(c echo.Context) error {
	return nil
}

func (s *Socket) HTTPForward(c echo.Context) error {
	id := c.Param("id")

	channel, ok := s.channels[id]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	msgBody := ""
	err := c.Bind(&msgBody)
	if err != nil {
		return err
	}

	msgCh := make(chan *Message)
	doneCh := make(chan bool)

	h := hashids.NewWithData(msgHd)
	msgMaxId++
  e, err := h.EncodeInt64([]int64{msgMaxId})
	if err != nil {
		return err
	}

	go channel.Connect(msgCh, doneCh)

	msgCh <- &Message{To: id, From: e, Body: msgBody}
	response := <-msgCh
	doneCh <- true

	return c.JSON(200, response)
}
