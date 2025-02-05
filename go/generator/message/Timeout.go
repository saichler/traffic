package message

import (
	"github.com/saichler/shared/go/share/interfaces"
	"time"
)

func (this *Message) StartTimeout(timeout int, log interfaces.ILogger) {
	log.Info("Setting Timeout to ", timeout, " seconds")
	time.Sleep(time.Second * time.Duration(timeout))
	if !this.complete {
		log.Warning("Reached timeout!")
		this.timeoutReached = true
		this.cond.L.Lock()
		defer this.cond.L.Unlock()
		log.Info("Broadcasting ....")
		this.cond.Broadcast()
		this.cond.Broadcast()
	}
}
