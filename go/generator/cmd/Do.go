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
}

func (this *Do) Name() string {
	return "Do"
}
func (this *Do) Help() string {
	return "Do a command"
}
func (this *Do) Run(log interfaces.ILogger) {
	if this.Udp_port == 0 {
		log.Error("Udp_port_local cannot be zero")
		return
	}
	if this.Destination == "" {
		log.Error("Destination cannot be blank")
		return
	}
	if this.Port == 0 {
		log.Error("Port cannot be blank or zero")
		return
	}
	if this.Quantity == 0 {
		log.Error("Quantity cannot be blank or zero")
		return
	}

	timeout := this.Quantity / 5000
	if timeout < 3 {
		timeout = 3
	}

	msg := message.NewCommand(message.Execute, this.Destination, this.Port, this.Quantity, "", timeout)
	Udp, err := udp.New(this.Udp_port-1, log, msg, true)
	if err != nil {
		log.Error(err.Error())
		return
	}

	msg.Lock()
	go msg.StartTimeout(timeout+1, log)
	err = Udp.Send(msg.String(), "127.0.0.1", this.Udp_port)
	if err != nil {
		log.Error(err.Error())
		return
	}
	msg.Wait(log)
	log.Info(msg.Text())
}
