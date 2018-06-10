package addon_gocui

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client defines a client that communicates with a Phoenix server over a channel/websocket
type Client struct {
	ch    chan *Payload
	conn  *websocket.Conn
	room  string
	topic string
}

// Message defines the structure of the messages sent across the channel
type Message struct {
	// Topic defines the Phoenix topic we're communicating on (must be 'phoenix' when sending heartbeats)
	Topic string `json:"topic"`

	// Event defines the type of event we're sending - one of phx_join, heartbeat, or chat_msg
	Event string `json:"event"`

	// Payload defines the body of the message coming across
	Payload Payload `json:"payload"`

	// Ref is a unique string that phoenix uses
	Ref string `json:"ref"`
}

// Payload defines the non-channel-specific body of the Message
type Payload struct {
	Status   string `json:"status"`
	Body     string `json:"body"`
	Username string `json:"username"`
}

// NewClient creates a Client that connects to a phoenix server over a given websocket address
func NewClient(address, room string) (*Client, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		return nil, err
	}

	client := Client{
		ch:   make(chan *Payload),
		conn: conn,
		room: room,
	}

	return &client, nil
}

// Listen listens for messages coming across the websocket. It exposes new chat messages via
// the client's channel
func (c *Client) Listen() error {
	err := c.join(c.room)
	if err != nil {
		return err
	}
	for {
		time.Sleep(time.Millisecond * 100)
		message, err := c.read()
		if err != nil {
			logrus.Warnf("Error reading from socket: %v", err)
		}
		if message != nil && message.Event == "chat_msg" {
			logrus.Debug("Received message from socket, sending to UI")
			c.ch <- &message.Payload
		}
	}
}

// Channel gives access to the channel holding new messages so that the UI can display them
func (c *Client) Channel() chan *Payload {
	return c.ch
}

// Join joins a given chat room via a join Message on the channel
func (c *Client) join(room string) error {
	topic := "rooms:" + room
	m := Message{
		Topic:   topic,
		Event:   "phx_join",
		Payload: Payload{},
	}
	err := c.write(&m)
	if err != nil {
		return err
	}
	logrus.Infof("Joined chat room %s", room)
	c.topic = topic
	go c.sendHeartbeats()
	return nil
}

// Read reads a Message from the channel
func (c *Client) read() (*Message, error) {
	var msg Message
	err := c.conn.ReadJSON(&msg)
	return &msg, err
}

func (c *Client) sendHeartbeats() {
	m := Message{
		Topic: "phoenix",
		Event: "heartbeat",
	}

	t := time.NewTicker(time.Second * 5)
	go func() {
		for range t.C {
			logrus.Debug("Sending heartbeat")
			c.write(&m)
		}
	}()
}

// SendMessage sends a new message to the server
func (c *Client) SendMessage(message string) error {
	p := Payload{
		Body: message,
	}
	m := Message{
		Topic:   c.topic,
		Event:   "chat_msg",
		Payload: p,
	}
	logrus.Debugf("Sending message '%s'", message)
	return c.write(&m)
}

func (c *Client) write(m *Message) error {
	return c.conn.WriteJSON(m)
}
