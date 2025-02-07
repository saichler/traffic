package cmd

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/traffic/go/generator/message"
	"github.com/saichler/traffic/go/generator/tcp"
	"github.com/saichler/traffic/go/generator/udp"
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
func (this *Start) Run(log interfaces.ILogger) string {

	if this.Udp_port <= 1000 {
		return log.Error("Udp_port must be less than 1000").Error()
	}
	if this.Tcp_port <= 1000 {
		return log.Error("Tcp_port must be less than 1000").Error()
	}
	msg := message.NewRequest(0, "Service")
	_, err := udp.New(this.Udp_port, log, msg, false)
	if err != nil {
		return log.Error(err).Error()
	}

	tcp.RunTcpServer(this.Tcp_port, log)
	return ""
}
