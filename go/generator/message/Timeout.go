package message

import (
	"github.com/saichler/shared/go/share/interfaces"
	"time"
)

func (this *Message) StartTimeout(timeout int, log interfaces.ILogger) {
	log.Info("Setting Timeout to ", timeout, " seconds")
	time.Sleep(time.Second * time.Duration(timeout))
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if !this.complete {
		log.Warning("Reached timeout!")
		this.timeoutReached = true
		log.Info("Broadcasting ....")
		this.cond.Broadcast()
	}
}
