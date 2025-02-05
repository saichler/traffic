package message

import (
	"github.com/saichler/shared/go/share/maps"
	"github.com/saichler/shared/go/share/strings"
	"net"
	"strconv"
	"sync"
	"time"
)

func (this *Message) Handle(packet []byte, addr *net.UDPAddr, protocol Protocol) {
	msg := MessageOf(string(packet), protocol, addr)
	if msg == nil {
		return
	}
	switch msg.action {
	case Execute:
		this.execute(msg, protocol, addr)
	case Request:
		this.request(msg, protocol, addr)
	case Response:
		this.response(msg, protocol, addr)
	default:
		panic(msg.String())
	}
}

func (this *Message) execute(msg *Message, protocol Protocol, addr *net.UDPAddr) {
	start := time.Now().Unix()
	curr := protocol.Port() + 1
	if this.cond == nil {
		this.cond = sync.NewCond(&sync.Mutex{})
	}

	this.pendingReply = maps.NewSyncMap()
	this.Lock()
	this.complete = false
	go this.StartTimeout(msg.timeout, protocol.Log())

	for i := 0; i < msg.quantity; i++ {
		var newProtocol Protocol
		for {
			newProtocol = protocol.New(curr, protocol.Log(), this, true)
			if newProtocol != nil {
				break
			}
			curr++
			if curr >= 65535 {
				protocol.Log().Warning("Deputed ports, sleeping for a second Milliseconds ", curr)
				time.Sleep(time.Millisecond * 250)
				curr = protocol.Port() + 1
			}
		}
		destMsg := NewRequest(i, newProtocol.String())
		this.pendingReply.Put(destMsg.id, newProtocol)
		err := newProtocol.Send(destMsg.String(), msg.destination, msg.port)
		if err != nil {
			Error(msg, err.Error(), protocol, addr)
		}
		curr++

		if i%5000 == 0 && i != 0 {
			protocol.Log().Info("Sleeping 250 Milliseconds Each 5K packets")
			time.Sleep(time.Millisecond * 250)
			curr = protocol.Port() + 1
		}
	}

	this.Wait(protocol.Log())
	this.complete = true
	end := time.Now().Unix()
	took := end - start
	responseMsg := strings.New("Total Sent: ", strconv.Itoa(msg.quantity), " Received:",
		strconv.Itoa(msg.quantity-this.pendingReply.Size()), " Took:", took, " Seconds.")
	resp := NewResponse(this.id, responseMsg.String())

	protocol.Reply(resp.String(), addr)

	all := this.pendingReply.Clean()
	for _, p := range all {
		proto, ok := p.(Protocol)
		if ok {
			proto.Shutdown()
		}
	}
}

func (this *Message) request(msg *Message, protocol Protocol, addr *net.UDPAddr) {
	//protocol.Log().Info(msg.String())
	resp := NewResponse(msg.id, protocol.String())
	protocol.Reply(resp.String(), addr)
}

func (this *Message) response(msg *Message, protocol Protocol, addr *net.UDPAddr) {
	msg.complete = true
	//protocol.Log().Info(msg.text)
	this.text = msg.text
	if protocol.Disposable() {
		protocol.Shutdown()
	}
	if this.pendingReply != nil {
		this.pendingReply.Delete(msg.id)
	}
	if this.pendingReply == nil || this.pendingReply.Size() == 0 {
		this.cond.Broadcast()
	}
}
