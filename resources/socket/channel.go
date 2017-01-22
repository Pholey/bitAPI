package socket

import (
  "io"
  "fmt"
  "golang.org/x/net/websocket"
  "github.com/speps/go-hashids"
)

const channelBufSize = 4096
var channelMaxId int64 = 0
var channelHd = &hashids.HashIDData{hashids.DefaultAlphabet, 64, "channel"}

type Channel struct {
  id         string
  outbox     map[string][]*Message
  recipients map[string]chan *Message
  ws         *websocket.Conn
  socket     *Socket
  ch         chan *Message
  doneCh     chan bool
}

func NewChannel(ws *websocket.Conn, socket *Socket) (*Channel, error) {
  if ws == nil {
    return nil, fmt.Errorf("No socket provided for new channel.")
  }

  outbox := make(map[string][]*Message)
  recipients := make(map[string]chan *Message)

  ch := make(chan *Message, channelBufSize)
  doneCh := make(chan bool)

  h := hashids.NewWithData(channelHd)
  channelMaxId++
  e, err := h.EncodeInt64([]int64{channelMaxId})

  if err != nil {
    return nil, err
  }

  return &Channel{e, outbox, recipients, ws, socket, ch, doneCh}, nil
}

func (c *Channel) HasRecipient(id string) bool {
  _, ok := c.recipients[id]
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

func (c *Channel) Connect(msgCh chan *Message, doneCh chan bool) error {
  for {
    select {
    case msg := <-msgCh:
      c.recipients[msg.From] = msgCh
      c.Write(msg)

      // queue, ok := c.outbox[msg.From]
      // if ok {
      //   for _, msg := range(queue) {
      //     msgCh <- msg
      //   }
      // }
    case <-doneCh:
      c.doneCh <- true
      return nil
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
      var msg = &Message{}
      err := websocket.JSON.Receive(c.ws, msg)
      if err == io.EOF {
        c.doneCh <- true
      } else if err != nil {
        c.socket.Err(err)
      } else if recipient, ok := c.recipients[msg.To]; ok {
        recipient <- msg
      } else if c.socket.HasRecipient(msg.To) {
        c.socket.send(msg)
      }
      // else {
      //   outBuf, ok := c.outbox[msg.To]
      //   if !ok {
      //     outBuf = make([]*Message)
      //   }
      //   c.outbox[msg.To] = append(outBuf, msg)
      // }
    }
  }
}
