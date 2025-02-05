package cmd

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/traffic/go/generator/message"
	"github.com/saichler/traffic/go/generator/udp"
	"time"
)

type Start struct {
	Udp_port int
}

func (this *Start) Name() string {
	return "Start"
}
func (this *Start) Help() string {
	return "start the service"
}
func (this *Start) Run(log interfaces.ILogger) {

	if this.Udp_port <= 1000 {
		log.Error("Udp_port must be less than 1000")
		return
	}
	msg := message.NewRequest(0, "Service")
	_, err := udp.New(this.Udp_port, log, msg, false)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
