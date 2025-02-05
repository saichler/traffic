package udp

import (
	"github.com/saichler/layer8/go/overlay/protocol"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/strings"
	"github.com/saichler/traffic/go/generator/message"
	"net"
	"strconv"
	"time"
)

type UDP struct {
	port       int
	addr       string
	conn       *net.UDPConn
	log        interfaces.ILogger
	shutdown   bool
	ml         message.IMessageListener
	disposable bool
}

func New(port int, log interfaces.ILogger, ml message.IMessageListener, disposable bool) (*UDP, error) {
	udp := &UDP{}
	udp.ml = ml
	udp.disposable = disposable
	if ml == nil {
		panic("")
	}
	udp.port = port
	udp.log = log
	err := udp.bind()
	if err != nil {
		return nil, err
	}
	go udp.rx()
	udp.addr = protocol.MachineIP
	return udp, nil
}

func (this *UDP) New(port int, log interfaces.ILogger, ml message.IMessageListener, disposable bool) message.Protocol {
	udp, err := New(port, log, ml, disposable)
	if err != nil {
		return nil
	}
	return udp
}

func (this *UDP) bind() error {
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(this.port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	this.conn = conn
	return nil
}

func (this *UDP) rx() {
	packet := make([]byte, 1024)
	defer this.conn.Close()
	if !this.disposable {
		this.log.Info("Starting UDP listener on port " + strconv.Itoa(this.port))
	}
	for !this.shutdown {
		n, addr, err := this.conn.ReadFromUDP(packet)
		if this.shutdown {
			break
		}
		if err != nil {
			this.log.Error(err.Error())
			break
		}
		if n > 0 {
			this.ml.Handle(packet[:n], addr, this)
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
	if !this.disposable {
		this.log.Info("End UDP listener on port " + strconv.Itoa(this.port))
	}
}

func (this *UDP) Log() interfaces.ILogger {
	return this.log
}

func (this *UDP) Send(msg, host string, port int) error {
	if this.ml == nil {
		panic("")
	}
	addr, err := net.ResolveUDPAddr("udp", host+":"+strconv.Itoa(port))
	if err != nil {
		return this.log.Error(err.Error())
	}
	this.Reply(msg, addr)
	return nil
}

func (this *UDP) Reply(msg string, addr *net.UDPAddr) {
	n, err := this.conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		this.log.Error(err.Error())
		return
	}
	if n != len(msg) {
		this.log.Error(strings.New("Number of bytes written is "+strconv.Itoa(n), " vs. ",
			strconv.Itoa(len(msg))).String())
	}
}

func (this *UDP) Shutdown() {
	if !this.disposable {
		this.log.Info("Shutting down UDP listener ", this.String())
	}
	this.shutdown = true
	this.conn.Close()
}

func (this *UDP) Port() int {
	return this.port
}

func (this *UDP) String() string {
	return strings.New(this.addr, ";", strconv.Itoa(this.port)).String()
}

func (this *UDP) Disposable() bool {
	return this.disposable
}
