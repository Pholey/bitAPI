package socket

type Message struct {
  To   string `json:"to"`
  From string `json:"from"`
  Body string `json:"body"`
}

func (m *Message) String() string {
  return "Message to " + m.To + " from " + m.From + " containing: " + m.Body
}
