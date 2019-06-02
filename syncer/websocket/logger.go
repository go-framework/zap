package websocket

import (
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Name.
	Name = "websocket"

	// Time allowed to write a envelope to the peer.
	WriteWait = 10 * time.Second
	// // Timeout for waiting on pong.
	PongWait = 60 * time.Second
	// Send pings to peer with this period.
	PingPeriod = (60 * time.Second * 9) / 10
	// Maximum envelope size allowed from peer.
	MaxMessageSize = 512
	// The max amount of messages.
	MessageBufferSize = 256
)

// envelope.
type envelope struct {
	t    int    // envelope type.
	data []byte // envelope data.
}

// Connect handler.
type ConnectHandler func()

// Redefined lumberjack Logger.
type Logger struct {
	// websocket conn.
	conn *websocket.Conn
	open bool
	// reconnect status.
	reconnect bool
	// have visitor.
	haveVisitor bool
	exit        chan struct{}
	output      chan *envelope
	waiting     chan struct{}
	rwMutex     *sync.RWMutex

	// connect handler.
	connectHandler ConnectHandler

	// Websocket server url.
	Url string `json:"url" yaml:"url" mapstructure:"url"`
	// Websocket write when have visitor, default is true.
	WriteOnVisitor bool `json:"write_on_visitor" yaml:"write_on_visitor" mapstructure:"write_on_visitor"`
	// Time allowed to write a envelope to the peer.
	WriteWait time.Duration `json:"write_wait" yaml:"write_wait" mapstructure:"write_wait"`
	// // Timeout for waiting on pong.
	PongWait time.Duration `json:"pong_wait" yaml:"pong_wait" mapstructure:"pong_wait"`
	// Send pings to peer with this period.
	PingPeriod time.Duration `json:"ping_period" yaml:"ping_period" mapstructure:"ping_period"`
	// Maximum envelope size allowed from peer.
	MaxMessageSize int64 `json:"max_message_size" yaml:"max_message_size" mapstructure:"max_message_size"`
	// The max amount of messages.
	MessageBufferSize int `json:"message_buffer_size" yaml:"message_buffer_size" mapstructure:"message_buffer_size"`
}

// New logger.
func New(url string) *Logger {
	l := &Logger{
		Url:               url,
		WriteOnVisitor:    true,
		WriteWait:         WriteWait,
		PongWait:          PongWait,
		PingPeriod:        PingPeriod,
		MaxMessageSize:    MaxMessageSize,
		MessageBufferSize: MessageBufferSize,
		output:            make(chan *envelope, MessageBufferSize),
		rwMutex:           &sync.RWMutex{},
		exit:              make(chan struct{}),
	}

	// go start connect.
	go l.connect()

	return l
}

// Get default logger.
func GetDefault() *Logger {
	l := &Logger{
		WriteWait:         WriteWait,
		WriteOnVisitor:    true,
		PongWait:          PongWait,
		PingPeriod:        PingPeriod,
		MaxMessageSize:    MaxMessageSize,
		MessageBufferSize: MessageBufferSize,
		output:            make(chan *envelope, MessageBufferSize),
		rwMutex:           &sync.RWMutex{},
		exit:              make(chan struct{}),
	}

	return l
}

// Implement Cloner interface.
func (l *Logger) Clone() io.Writer {
	n := *l

	n.output = make(chan *envelope, l.MessageBufferSize)
	n.rwMutex = &sync.RWMutex{}
	n.exit = make(chan struct{})

	// go start connect.
	go n.connect()

	return &n
}

// Waiting complete.
func (l *Logger) Waiting() {
	// waiting complete.
	if l.waiting == nil {
		l.waiting = make(chan struct{})

	}
	<-l.waiting
	close(l.waiting)
	l.waiting = nil
}

// Implement Writer interface.
func (l *Logger) Write(p []byte) (n int, err error) {
	n = len(p)

	// write enable?
	if !l.haveVisitor && l.WriteOnVisitor {
		return
	}

	// copy slice.
	b := make([]byte, n)
	copy(b, p)

	// write envelope to channel.
	err = l.writeMessage(&envelope{t: websocket.TextMessage, data: b})
	if err != nil {
		return
	}

	return
}

// Close.
func (l *Logger) Close() error {
	if l.closed() {
		return errors.New("websocket conn is already closed")
	}

	defer func() {
		_ = l.write(&envelope{t: websocket.CloseMessage, data: nil})
		l.close()
	}()

	// waiting send complete.
	if l.waiting == nil {
		l.waiting = make(chan struct{})
	}
	<-l.waiting
	close(l.waiting)
	l.waiting = nil

	return nil
}

// Set url.
func (l *Logger) SetUrl(url string) {
	l.Url = url
}

// Set connect handler.
func (l *Logger) SetConnectHandler(handler ConnectHandler) {
	l.connectHandler = handler
}

// connect to server.
func (l *Logger) connect() {

	d := 1 * time.Second
	ticker := time.NewTicker(d)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case _, ok := <-l.exit:
			if !ok {
				return
			}
		case <-ticker.C:
			if err := l.dial(); err != nil {
				d *= 2
				if d > WriteWait {
					d = 1 * time.Second
				}
				ticker = time.NewTicker(d)
				break
			}

			l.open = true

			// connect handler callback.
			if l.connectHandler != nil {
				l.connectHandler()
			}

			// send waiting signal.
			if l.waiting != nil {
				l.waiting <- struct{}{}
			}

			return
		}
	}
}

// dial to server.
func (l *Logger) dial() (err error) {
	// get websocket dialer
	mDial := websocket.DefaultDialer
	// set write buffer size
	mDial.WriteBufferSize = int(l.MaxMessageSize)

	// dial get conn
	l.conn, _, err = mDial.Dial(l.Url+"?role=sender", http.Header{
		"role": []string{"sender"},
	})
	if err != nil {
		return
	}

	// read pump.
	go l.readPump()

	// write pump.
	go l.writePump()

	return
}

// write envelope to channel.
func (l *Logger) writeMessage(message *envelope) error {

	select {
	case l.output <- message:
	default:
		return errors.New("envelope buffer is full")
	}

	return nil
}

// write envelope.
func (l *Logger) write(message *envelope) error {
	if l.closed() {
		return errors.New("tried to write to a closed websocket conn")
	}

	_ = l.conn.SetWriteDeadline(time.Now().Add(l.WriteWait))

	return l.conn.WriteMessage(message.t, message.data)
}

// readPump pumps messages from the websocket connection.
func (l *Logger) readPump() {
	defer func() {
		// restart.
		if l.reconnect {
			l.restart()
		}
	}()

	l.conn.SetReadLimit(l.MaxMessageSize)
	_ = l.conn.SetReadDeadline(time.Now().Add(l.PongWait))

	l.conn.SetPongHandler(func(string) error {
		_ = l.conn.SetReadDeadline(time.Now().Add(l.PongWait))
		return nil
	})

	for {
		select {
		case _, ok := <-l.exit:
			if !ok {
				return
			}
		default:
		}
		// read envelope.
		_, data, err := l.conn.ReadMessage()
		if err != nil {
			l.reconnect = true
			return
		}
		// new message.
		msg, err := newMessage(data)
		if err != nil {
			continue
		}
		// 命令
		if msg.Type == 0 {
			// 100 - VisitorOnline
			// 101 - NoVisitor
			switch msg.Command {
			case 100:
				l.haveVisitor = true
			case 101:
				l.haveVisitor = false
			}
		}

	}
}

// writePump pumps messages to the websocket connection.
func (l *Logger) writePump() {
	ticker := time.NewTicker(l.PingPeriod)

	defer func() {
		ticker.Stop()
		// restart.
		if l.reconnect {
			l.restart()
		}
	}()

	for {
		select {
		case _, ok := <-l.exit:
			if !ok {
				return
			}
		case <-ticker.C:
			// ping.
			if err := l.write(&envelope{t: websocket.PingMessage, data: nil}); err != nil {
				l.reconnect = true
				return
			}
		case message, ok := <-l.output:
			if !ok {
				return
			}
			// write raw envelope.
			if err := l.write(message); err != nil {
				l.reconnect = true
				return
			}
			// waiting
			if l.waiting != nil && len(l.output) == 0 {
				l.waiting <- struct{}{}
			}
		}
	}
}

// is closed?
func (l *Logger) closed() bool {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()

	return !l.open
}

// close.
func (l *Logger) close() {
	if !l.closed() {
		l.rwMutex.Lock()
		l.open = false
		close(l.exit)
		_ = l.conn.Close()
		close(l.output)
		l.rwMutex.Unlock()
	}
}

// restart.
func (l *Logger) restart() {
	if !l.closed() {
		l.rwMutex.Lock()
		l.open = false
		close(l.exit)
		_ = l.conn.Close()
		l.exit = make(chan struct{})
		go l.connect()
		l.rwMutex.Unlock()
	}
}
