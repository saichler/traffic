package tests

import (
	"testing"
)

func TestTcp_single_packet(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Tcp_port=" + TCP_Port_1, "Destination=127.0.0.1", "Port=" + TCP_Port_2, "Quantity=1"}
	if !testCMD(args, "Total TCP Sent:1 OK:1 Err:0 Timeout:false", t) {
		return
	}
}

func Test100PacketTCP(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Tcp_port=" + TCP_Port_1, "Destination=127.0.0.1", "Port=" + TCP_Port_2, "Quantity=100"}
	if !testCMD(args, "Total TCP Sent:100 OK:100 Err:0 Timeout:false", t) {
		return
	}
}

func Test100PacketTCPErr(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Tcp_port=" + TCP_Port_1, "Destination=127.0.0.2", "Port=" + TCP_Port_2, "Quantity=100"}
	if !testCMD(args, "Total TCP Sent:100 OK:0 Err:100 Timeout:false", t) {
		return
	}
}

func Test1000PacketTCP(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Tcp_port=" + TCP_Port_1, "Destination=127.0.0.1", "Port=" + TCP_Port_2, "Quantity=1000"}
	if !testCMD(args, "Total TCP Sent:1000 OK:1000 Err:0 Timeout:false", t) {
		return
	}
}

func Test10000PacketTCP(t *testing.T) {
	args := []string{"generator", "Do", "Udp_port=" + UDP_Port_1, "Tcp_port=" + TCP_Port_1, "Destination=127.0.0.1", "Port=" + TCP_Port_2, "Quantity=10000"}
	if !testCMD(args, "Total TCP Sent:10000 OK:10000 Err:0 Timeout:false", t) {
		return
	}
}
