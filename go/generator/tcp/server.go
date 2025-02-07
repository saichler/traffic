package tcp

import (
	"bytes"
	"github.com/saichler/shared/go/share/interfaces"
	"net/http"
	"strconv"
)

type TcpSerer struct {
	port      int
	webServer *http.Server
}

func init() {
	http.DefaultServeMux.HandleFunc("/", response)
}

func RunTcpServer(port int, log interfaces.ILogger) string {
	server := &TcpSerer{}
	server.webServer = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: http.DefaultServeMux,
	}
	log.Info("Starting TCP listener on port " + strconv.Itoa(port))
	err := server.webServer.ListenAndServe()
	if err != nil {
		log.Error("Failed TCP listener on port " + strconv.Itoa(port))
		return err.Error()
	}
	return ""
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
