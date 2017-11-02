package socket

import (
  "io"
  "fmt"
  "golang.org/x/net/websocket"
  db "github.com/Pholey/bitAPI/db"
)

const channelBufSize = 4096

type Channel struct {
  id         string
  recipients map[string]chan *Message
  ws         *websocket.Conn
  socket     *Socket
  ch         chan *Message
  doneCh     chan bool
}

func NewChannel(ws *websocket.Conn, socket *Socket, userId int64, channelName string) (*Channel, error) {
  if ws == nil {
    return nil, fmt.Errorf("No socket provided for new channel.")
  }

  _, err := db.Session.
    InsertBySql("INSERT INTO user_channel (id, channel_name) VALUES (?, ?) ON CONFLICT (channel_name) DO NOTHING", userId, channelName).
    Exec()

  if (err != nil) {
    return nil, err
  }

  recipients := make(map[string]chan *Message)

  ch := make(chan *Message, channelBufSize)
  doneCh := make(chan bool)

  return &Channel{channelName, recipients, ws, socket, ch, doneCh}, nil
}

func (c *Channel) HasRecipient(id string) bool {
  _, ok := c.recipients[id]
  fmt.Println("Really? ", ok)
  return ok
}

func (c *Channel) Write(msg *Message) error {
  select {
  case c.ch <- msg:
  default:
    c.socket.Del(c)
    err := fmt.Errorf("channel %d is disconnected.", c.id)
    return err
  }
  return nil
}

func (c *Channel) Connect(msgCh chan *Message) error {
  for {
    select {
    case msg := <-msgCh:
      c.recipients[msg.From] = msgCh
      c.Write(msg)
    }
  }
}

func (c *Channel) Done() {
  c.doneCh <- true
}

func (c *Channel) Listen() {
  go c.listenWrite()
  c.listenRead()
}

func (c *Channel) listenWrite() {
  for {
    select {
    case msg := <-c.ch:
      websocket.JSON.Send(c.ws, msg)
    case <-c.doneCh:
      c.socket.Del(c)
      c.doneCh <- true
      return
    }
  }
}

func (c *Channel) listenRead() {
  for {
    select {
    case <-c.doneCh:
      c.socket.Del(c)
      c.doneCh <- true
      return
    default:

      var msg = &Message{Body: &SDP{}}
      err := websocket.JSON.Receive(c.ws, msg)
      if err == io.EOF {
        c.doneCh <- true
      } else if err != nil {
        c.socket.Err(err)
      } else if recipient, ok := c.recipients[msg.To]; ok {
        recipient <- msg
        delete(c.recipients, msg.To)
      } else if c.socket.HasRecipient(msg.To) {
        c.socket.send(msg)
      }
    }
  }
}
