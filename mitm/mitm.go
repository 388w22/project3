//go:build linux

// Take a close look around the starter code,
// and feel free to take a deep dive into tests.
// Understanding what the given unit tests do
// will be a huge help in pulling off your attack.

package main

// These are the imports we used, but feel free to use anything from
// gopacket or the Go standard libraries. DO NOT import other third-party
// libraries, as your code may fail to compile on the autograder.
import (
	"net"
	"net/http"
	"os"

	"bank.com/mitm/network"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/sys/unix"
)

// ==============================
//  DNS MITM PORTION
// ==============================

// startDNSServer begins listening to traffic on the
// main network interface, handing off any DNS packets
// it finds to handleDNSPacket.
//
// You do not need to modify this function.
func startDNSServer() {
	handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	err = handle.SetBPFFilter("udp") // only grab UDP packets
	if err != nil {
		// More on BPF filtering:
		// https://www.ibm.com/support/knowledgecenter/SS42VS_7.4.0/com.ibm.qradar.doc/c_forensics_bpf.html
		panic(err)
	}
	defer handle.Close()

	// Loop over each packet received
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for pkt := range packetSource.Packets() {
		dns := pkt.Layer(layers.LayerTypeDNS)
		if dns != nil {
			handleDNSPacket(pkt)
		}
	}
}

// handleDNSPacket is called each time a DNS packet is observed
// on the network (both inbound and outbound).
func handleDNSPacket(packet gopacket.Packet) {
	panic("handleDNSPacket unimplemented!")

	// See dns.go for some possibly-relevant functions
	// (some of which you'll need to implement!).

	// Use sendRawUDP to send the final
	// response packet to its destination.
}

// sendRawUDP sends the data of an IP packet specified by data
// to the IP address and UDP port specified by ip and port.
//
// You do not need to modify this function.
func sendRawUDP(ip net.IP, port layers.UDPPort, data []byte) {
	// Opens an IPv4 socket to destination host/port.
	sock, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)
	if err != nil {
		panic(err)
	}
	addr := unix.SockaddrInet4{
		Port: int(port),
		Addr: *(*[4]byte)(ip.To4()),
	}
	if err := unix.Sendto(sock, data, 0, &addr); err != nil {
		panic(err)
	}
	if err := unix.Close(sock); err != nil {
		panic(err)
	}
}

// ==============================
//  HTTP MITM PORTION
// ==============================

// startHTTPServer sets up and hosts a basic HTTP server
// which calls handleHTTP for each request.
//
// You do not need to modify this function.
func startHTTPServer() {
	http.HandleFunc("/", handleHTTP)
	panic(http.ListenAndServe(":80", nil))
}

// handleHTTP is called each time a request is made
// to the local HTTP server.
//
// (Remember: you are pretending to be bank.com.)
func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/kill" {
		os.Exit(1)
	}

	bank := "http://" + network.GetBankIP().String()

	// See http.go for some possibly-relevant functions
	// (which you'll need to implement!).

	panic("handleHTTP unimplemented!")
}

func main() {
	// The DNS server is run concurrently alongside
	// the HTTP server as a goroutine
	go startDNSServer()

	startHTTPServer()
}
