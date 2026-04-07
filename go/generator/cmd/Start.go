package cmd

import (
	"github.com/saichler/l8traffic/go/generator/message"
	"github.com/saichler/l8traffic/go/generator/tcp"
	"github.com/saichler/l8traffic/go/generator/udp"
	"github.com/saichler/l8types/go/ifs"
	"os"
	"os/signal"
	"syscall"
)

type Start struct {
	Udp_port int
	Tcp_port int
}

func (this *Start) Name() string {
	return "Start"
}
func (this *Start) Help() string {
	return "start the service"
}
func (this *Start) Run(log ifs.ILogger) string {

	if this.Udp_port <= 1000 {
		return log.Error("Udp_port must be greater than 1000").Error()
	}
	if this.Tcp_port <= 1000 {
		return log.Error("Tcp_port must be greater than 1000").Error()
	}
	msg := message.NewRequest(0, "Service")
	udpConn, err := udp.New(this.Udp_port, log, msg, false)
	if err != nil {
		return log.Error(err).Error()
	}

	tcpServer := tcp.RunTcpServer(this.Tcp_port, log)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info("Shutting down...")
	tcpServer.Shutdown()
	udpConn.Shutdown()
	return ""
}
