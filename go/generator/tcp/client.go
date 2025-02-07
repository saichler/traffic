package tcp

import (
	"bytes"
	strings2 "github.com/saichler/shared/go/share/strings"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Workers struct {
	workers    int
	maxWorkers int
	responses  []string
	errors     []string
	okCount    int
	errCount   int
	took       int64
	cond       *sync.Cond
}

func newWorkers(w int) *Workers {
	workers := &Workers{}
	workers.cond = sync.NewCond(&sync.Mutex{})
	workers.workers = w
	workers.maxWorkers = w
	workers.responses = make([]string, 0)
	workers.errors = make([]string, 0)
	return workers
}

func SendHttpRequest(host string, port, quantity int) string {
	workers := newWorkers(100)
	start := time.Now().Unix()
	for i := 0; i < quantity; i++ {
		go workers.sendHttpRequest(host, port)
		workers.Wait()
	}
	workers.WaitEnd()
	end := time.Now().Unix()
	return CreateReport("TCP", quantity,
		workers.okCount, workers.errCount,
		int(end-start), workers.responses,
		workers.errors, false)
}

func (this *Workers) WaitEnd() {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	for this.workers < this.maxWorkers {
		this.cond.Wait()
	}
}

func (this *Workers) Wait() {
	this.cond.L.Lock()
	defer this.cond.L.Unlock()
	if this.workers == 0 {
		this.cond.Wait()
	}
	this.workers--
}

func (this *Workers) sendHttpRequest(host string, port int) {
	var resp string
	var conn net.Conn
	var err error

	defer func() {
		if conn != nil {
			conn.Close()
		}
		this.cond.L.Lock()
		defer this.cond.L.Unlock()
		if err != nil {
			this.errCount++
			if this.errCount <= 5 {
				this.errors = append(this.errors, err.Error())
			}
		} else {
			this.okCount++
			if this.okCount <= 5 {
				this.responses = append(this.responses, resp)
			}
		}

		this.cond.Broadcast()
		this.workers++
	}()

	request := "GET / HTTP/1.1\r\n" +
		"Host: www.example.com\r\n" +
		"User-Agent: MyCustomClient\r\n" +
		"Connection: close\r\n\r\n"
	dialer := net.Dialer{Timeout: time.Second * 2}
	conn, err = dialer.Dial("tcp", hostPort(host, port))
	if err != nil {
		return
	}

	_, err = conn.Write([]byte(request))
	if err != nil {
		return
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return
	}
	fullResp := string(buff[:n])
	index := strings.LastIndex(fullResp, "\n")
	resp = fullResp[index+1:]
}

func hostPort(host string, port int) string {
	buff := bytes.Buffer{}
	buff.WriteString(host)
	buff.WriteString(":")
	buff.WriteString(strconv.Itoa(port))
	return buff.String()
}

func CreateReport(proto string, sent, ok, err, took int, oks, errs []string, timeout bool) string {
	str := strings2.New()
	tmout := "false"
	if timeout {
		tmout = "true"
	}
	str.Add("Total ", proto, " Sent:", strconv.Itoa(sent),
		" OK:", strconv.Itoa(ok),
		" Err:", strconv.Itoa(err),
		" Timeout:", tmout,
		" Took:", strconv.Itoa(took), " seconds\n")
	str.Add("OK Sample:\n")
	for _, resp := range oks {
		str.Add(" - ")
		str.Add(resp)
		str.Add("\n")
	}
	str.Add("Err Sample:\n")
	for _, e := range errs {
		str.Add(" - ")
		str.Add(e)
		str.Add("\n")
	}
	return str.String()
}
