package message

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/maps"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	SendUDP  = "sendUdp"
	SendTCP  = "sendTcp"
	Request  = "request"
	Response = "response"
)

type Message struct {
	id             int
	action         string
	destination    string
	port           int
	quantity       int
	text           string
	timeout        int
	complete       bool
	timeoutReached bool
	timeoutGen     int

	pendingReply *maps.SyncMap
	cond         *sync.Cond
}

type IMessageListener interface {
	Handle([]byte, *net.UDPAddr, Protocol)
}

type Protocol interface {
	Reply(string, *net.UDPAddr)
	Send(string, string, int) error
	New(int, ifs.ILogger, IMessageListener, bool) Protocol
	Log() ifs.ILogger
	String() string
	Port() int
	Shutdown()
	Disposable() bool
}

// nextId is a global counter shared by both client and service messages.
// IDs only need to be unique within a request-response pair, so a shared
// counter is acceptable.
var nextId = 0
var mtx = &sync.Mutex{}

func nextID() int {
	mtx.Lock()
	defer mtx.Unlock()
	nextId++
	return nextId
}

func NewCommand(action string,
	destination string,
	port int,
	quantity int,
	text string,
	timeout int) *Message {
	msg := &Message{}
	msg.cond = sync.NewCond(&sync.Mutex{})
	msg.id = nextID()
	msg.action = action
	msg.destination = destination
	msg.port = port
	msg.quantity = quantity
	msg.text = text
	msg.timeout = timeout
	return msg
}

func NewRequest(i int, from string) *Message {
	msg := &Message{}
	msg.id = nextID()
	msg.action = Request
	msg.text = strings2.New("Request - ", strconv.Itoa(i), " from ", from).String()
	return msg
}

func NewResponse(id int, from string) *Message {
	msg := &Message{}
	msg.id = id
	msg.action = Response
	msg.text = strings2.New("Response - ", strconv.Itoa(id), " from ", from).String()
	return msg
}

func Error(msg *Message, errMsg string, p Protocol, addr *net.UDPAddr) {
	msg.text = errMsg
	msg.action = Response
	p.Reply(msg.String(), addr)
}

func MessageOf(str string, p Protocol, addr *net.UDPAddr) *Message {
	msg := &Message{}
	split := strings.SplitN(str, "|", 7)
	if len(split) < 7 {
		msg.id, _ = strconv.Atoi(split[0])
		Error(msg, "message format error:"+str, p, addr)
		return nil
	}
	msg.id, _ = strconv.Atoi(split[0])
	msg.action = split[1]
	msg.destination = split[2]
	msg.port, _ = strconv.Atoi(split[3])
	msg.quantity, _ = strconv.Atoi(split[4])
	msg.text = split[5]
	msg.timeout, _ = strconv.Atoi(split[6])
	return msg
}

func (this *Message) Bytes() []byte {
	str := this.String()
	return []byte(str)
}

func (this *Message) Text() string {
	if this.cond != nil {
		this.cond.L.Lock()
		defer this.cond.L.Unlock()
	}
	return this.text
}

func (this *Message) String() string {
	str := strings2.New()
	str.Add(strconv.Itoa(this.id)).
		Add("|").
		Add(this.action).
		Add("|").
		Add(this.destination).
		Add("|").
		Add(strconv.Itoa(this.port)).
		Add("|").
		Add(strconv.Itoa(this.quantity)).
		Add("|").
		Add(this.text).
		Add("|").
		Add(strconv.Itoa(this.timeout))
	return str.String()
}

func (this *Message) Lock() {
	this.cond.L.Lock()
}
func (this *Message) Wait(log ifs.ILogger) {
	log.Info("Waiting for task to finish...")
	defer this.cond.L.Unlock()
	this.cond.Wait()
	if !this.timeoutReached {
		this.complete = true
		log.Info("Finished waiting!")
	}
}
