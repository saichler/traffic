package cmd

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/traffic/go/generator/message"
	"github.com/saichler/traffic/go/generator/udp"
)

type Do struct {
	Udp_port    int
	Destination string
	Port        int
	Quantity    int
	Timeout     int
}

func (this *Do) Name() string {
	return "Do"
}
func (this *Do) Help() string {
	return "Do a command"
}
func (this *Do) Run(log interfaces.ILogger) string {
	if this.Udp_port == 0 {
		return log.Error("Udp_port cannot be zero").Error()
	}
	if this.Destination == "" {
		return log.Error("Destination cannot be blank").Error()
	}
	if this.Port == 0 {
		return log.Error("Port cannot be blank or zero").Error()
	}
	if this.Quantity == 0 {
		return log.Error("Quantity cannot be blank or zero").Error()
	}

	timeout := this.Timeout

	if this.Timeout == 0 {
		timeout = this.Quantity / 5000
		if timeout < 3 {
			timeout = 3
		}
	}

	msg := message.NewCommand(message.Execute, this.Destination, this.Port, this.Quantity, "", timeout)
	Udp, err := udp.New(this.Udp_port-1, log, msg, true)
	if err != nil {
		return log.Error(err.Error()).Error()
	}

	msg.Lock()
	go msg.StartTimeout(timeout+1, log)
	err = Udp.Send(msg.String(), "127.0.0.1", this.Udp_port)
	if err != nil {
		return log.Error(err.Error()).Error()
	}
	msg.Wait(log)
	log.Info(msg.Text())
	return msg.Text()
}
