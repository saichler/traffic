package tcp

import (
	"bytes"
	"github.com/saichler/l8types/go/ifs"
	"net/http"
	"strconv"
)

type TcpServer struct {
	port      int
	webServer *http.Server
}

func RunTcpServer(port int, log ifs.ILogger) *TcpServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/", response)
	server := &TcpServer{port: port}
	server.webServer = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: mux,
	}
	log.Info("Starting TCP listener on port " + strconv.Itoa(port))
	go func() {
		err := server.webServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("Failed TCP listener on port " + strconv.Itoa(port))
		}
	}()
	return server
}

func (this *TcpServer) Shutdown() {
	this.webServer.Close()
}

func response(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write(format(r))
}

func format(r *http.Request) []byte {
	buff := bytes.Buffer{}
	buff.WriteString("TCP Response ")
	buff.WriteString(r.RemoteAddr)
	return buff.Bytes()
}
