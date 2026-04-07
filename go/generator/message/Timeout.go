package message

import (
	"github.com/saichler/l8types/go/ifs"
	"time"
)

func (this *Message) StartTimeout(timeout int, log ifs.ILogger) {
	// Capture the current timeout generation under the cond lock.
	// If sendUdp starts a new round before this goroutine fires,
	// timeoutGen will have been incremented and this goroutine
	// will be a no-op.
	this.cond.L.Lock()
	gen := this.timeoutGen
	this.cond.L.Unlock()

	log.Info("Setting Timeout to ", timeout, " seconds")
	time.Sleep(time.Second * time.Duration(timeout))
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if !this.complete && this.timeoutGen == gen {
		log.Warning("Reached timeout!")
		this.timeoutReached = true
		log.Info("Broadcasting ....")
		this.cond.Broadcast()
	}
}
