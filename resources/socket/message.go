package socket

type SDP struct {
	SDPType string `json:"type"`
	SDP     string `json:"sdp"`
}

type Message struct {
  To   string `json:"to"`
  From string `json:"from"`
  Body *SDP `json:"body"`
}

func (s *SDP) String() string {
  return "{type: " + s.SDPType + ", sdp: " + s.SDP +"}"
}

func (m *Message) String() string {
  return "Message to " + m.To + " from " + m.From + " containing: " + m.Body.String()
}
