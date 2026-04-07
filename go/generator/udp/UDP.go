package udp

import (
	"errors"
	"github.com/saichler/l8traffic/go/generator/message"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/ipsegment"
	"github.com/saichler/l8utils/go/utils/strings"
	"net"
	"strconv"
	"sync/atomic"
)

type UDP struct {
	port       int
	addr       string
	conn       *net.UDPConn
	log        ifs.ILogger
	shutdown   atomic.Bool
	ml         message.IMessageListener
	disposable bool
}

func New(port int, log ifs.ILogger, ml message.IMessageListener, disposable bool) (*UDP, error) {
	udp := &UDP{}
	udp.ml = ml
	udp.disposable = disposable
	if ml == nil {
		return nil, errors.New("message Listener is required")
	}
	udp.port = port
	udp.log = log
	err := udp.bind()
	if err != nil {
		return nil, err
	}
	go udp.rx()
	udp.addr = ipsegment.MachineIP
	return udp, nil
}

func (this *UDP) New(port int, log ifs.ILogger, ml message.IMessageListener, disposable bool) message.Protocol {
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
	packet := make([]byte, 65535)
	defer this.conn.Close()
	if !this.disposable {
		this.log.Info("Starting UDP listener on port " + strconv.Itoa(this.port))
	}
	for !this.shutdown.Load() {
		n, addr, err := this.conn.ReadFromUDP(packet)
		if this.shutdown.Load() {
			break
		}
		if err != nil {
			this.log.Error(err.Error())
			break
		}
		if n > 0 {
			this.ml.Handle(packet[:n], addr, this)
		}
	}
	if !this.disposable {
		this.log.Info("End UDP listener on port " + strconv.Itoa(this.port))
	}
}

func (this *UDP) Log() ifs.ILogger {
	return this.log
}

func (this *UDP) Send(msg, host string, port int) error {
	if this.ml == nil {
		return errors.New("message Listener is required")
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
	this.shutdown.Store(true)
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
