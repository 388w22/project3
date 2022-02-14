package main

import (
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ProduceIPPacket returns the bytes of an IP packet with the
// IPv4, UDP and DNS layer details specified in ip, udp, and dns.
// You do not need to modify this function.
func ProduceIPPacket(ip *layers.IPv4, udp *layers.UDP, dns *layers.DNS) []byte {
	// Set up some basic constants for the context of this function.
	ip.Version = 4
	ip.Protocol = layers.IPProtocolUDP

	// The checksum for the level 4 header (which includes UDP) depends on
	// what level 3 protocol encapsulates it; let UDP know it will be wrapped
	// inside IPv4.
	if err := udp.SetNetworkLayerForChecksum(ip); err != nil {
		log.Panic(err)
	}

	// Now we're ready to seal off and send the packet.
	// Serialization refers to "flattening" a packet's different layers into a
	// raw stream of bytes to be sent over the network.
	// Here, we want to automatically populate length and checksum fields with the correct values.
	serializeOpts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	buf := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buf, serializeOpts, ip, udp, dns); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

// HasQuestionForDomain should return whether the DNS packet
// represented by dns contains a question for domain.
//
// Hint: hover your cursor over layers.DNS for more
//       information; also try scrolling to the bottom
//       of that popup and view its online documentation!
func HasQuestionForDomain(dns *layers.DNS, domain string) bool {
	panic("HasQuestionForDomain unimplemented!")
}

// AnswerForQuestion should return an answer corresponding
// to question which points to the IP address ip.
func AnswerForQuestion(question layers.DNSQuestion, ip net.IP) layers.DNSResourceRecord {
	panic("AnswerForQuestion unimplemented!")
}
