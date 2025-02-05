package tests

import (
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/traffic/go/generator/cmd"
	"strings"
	"testing"
)

var log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
var cmds = cmd.NewCommands(log)
var UDP_Port_1 = "47000"
var UDP_Port_2 = "48000"

func init() {
	args := []string{"generator", "Start", "Udp_port=" + UDP_Port_1}
	go cmds.Run(args)
	args = []string{"generator", "Start", "Udp_port=" + UDP_Port_2}
	go cmds.Run(args)

}

func TestUnknownCommand(t *testing.T) {
	args := []string{"generator", "unknowncommand", "Udp_port=" + UDP_Port_1}
	msg := cmds.Run(args)
	if msg == "" {
		log.Fail(t, "Expected an error message")
		return
	}
	if !strings.Contains(msg, "no unknowncommand") {
		log.Fail(t, "Unexpected error:", msg)
		return
	}
}

func testCMD(args []string, expected string, t *testing.T) bool {
	msg := cmds.Run(args)
	if !strings.Contains(msg, expected) {
		log.Error("Received:", msg)
		log.Fail(t, "Expected '", expected, "' in command reply outcome")
		return false
	}
	return true
}

func TestStartInvalidPort(t *testing.T) {
	args := []string{"generator", "Start", "Udp_port=999"}
	if !testCMD(args, "Udp_port must be less than 1000", t) {
		return
	}
}

func TestPrintCommandHelp(t *testing.T) {
	args := []string{"generator", "Do"}
	if !testCMD(args, "Do a command", t) {
		return
	}
	args = []string{"generator", "Start"}
	if !testCMD(args, "start the service", t) {
		return
	}
}

func TestInvalidCommandArgument(t *testing.T) {
	args := []string{"generator", "Do", "test"}
	if !testCMD(args, "Invalid argument format: test", t) {
		return
	}
}

func TestUnknownCommandArgument(t *testing.T) {
	args := []string{"generator", "Do", "test=3"}
	if !testCMD(args, "Invalid argument name: test", t) {
		return
	}
}

func TestMissingCommandArgument(t *testing.T) {
	args := []string{"generator", "Do", "Port=15"}
	if !testCMD(args, "Udp_port cannot be zero", t) {
		return
	}
	args = []string{"generator", "Do", "Udp_port=10"}
	if !testCMD(args, "Destination cannot be blank", t) {
		return
	}
	args = []string{"generator", "Do", "Udp_port=10", "Destination=127.0.0.1"}
	if !testCMD(args, "Port cannot be blank", t) {
		return
	}
	args = []string{"generator", "Do", "Udp_port=10", "Destination=127.0.0.1", "Port=10"}
	if !testCMD(args, "Quantity cannot be blank", t) {
		return
	}
}

func TestUdp_single_packet(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Destination=127.0.0.1", "Port=" + UDP_Port_2, "Quantity=1"}
	if !testCMD(args, "Total Sent: 1 Received:1 Took", t) {
		return
	}
}

func TestUdp_timeout(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Destination=127.0.0.1", "Port=1000", "Quantity=1", "Timeout=1"}
	if !testCMD(args, "Timeout!", t) {
		return
	}
}

func TestUdp_1000_packets(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Destination=127.0.0.1", "Port=" + UDP_Port_2, "Quantity=1000"}
	if !testCMD(args, "Total Sent: 1000 Received:1000 Took", t) {
		return
	}
}

func TestUdp_10000_packets(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Destination=127.0.0.1", "Port=" + UDP_Port_2, "Quantity=10000"}
	if !testCMD(args, "Total Sent: 10000 Received:10000 Took", t) {
		return
	}
}

/*
func TestUdp_100000_packets(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Destination=127.0.0.1", "Port=" + UDP_Port_2, "Quantity=100000"}
	if !testCMD(args, "Total Sent: 100000 Received:100000 Took", t) {
		return
	}
}*/
