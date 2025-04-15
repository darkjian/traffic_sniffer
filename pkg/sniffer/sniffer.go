package sniffer

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/songgao/water"
	"log"
)

type Sniffer struct {
	devType    water.DeviceType
	bufferSize int
}

func New(bufSize int) *Sniffer {
	return &Sniffer{devType: water.TUN, bufferSize: bufSize}
}

func (s *Sniffer) Run() {
	if s.bufferSize == 0 {
		s.bufferSize = 2048
	}

	// creating tun interface
	iFace, err := water.New(water.Config{DeviceType: s.devType})
	if err != nil {
		log.Fatalf("Error creating utun interface: %v", err)
	}

	log.Printf("Created interface: %s", iFace.Name())

	var (
		buffer = make([]byte, s.bufferSize)
		n      int // number of bytes
		layer  gopacket.Layer
	)

	for {
		if n, err = iFace.Read(buffer); err != nil {
			log.Fatalf("Error reading from TUN interface: %v", err)
		}

		packet := gopacket.NewPacket(
			buffer[:n],
			layers.LayerTypeIPv4,
			gopacket.Default,
		)

		if layer = packet.Layer(layers.LayerTypeIPv4); layer == nil {
			// skip if not IPv4
			continue
		}

		ip, _ := layer.(*layers.IPv4)

		switch ip.Protocol {
		case layers.IPProtocolTCP:
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)
				log.Printf(
					"[TCP] IP: %s -> %s, Port: %s -> %s, Bytes: %d, Raw: %v",
					ip.SrcIP.String(), ip.DstIP.String(), tcp.SrcPort.String(), tcp.DstPort.String(), n, buffer[:n],
				)
			}
		case layers.IPProtocolUDP:
			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				udp, _ := udpLayer.(*layers.UDP)
				log.Printf(
					"[UDP] IP: %s -> %s, Port: %s -> %s, Bytes: %d, Raw: %v",
					ip.SrcIP.String(), ip.DstIP.String(), udp.SrcPort.String(), udp.DstPort.String(), n, buffer[:n],
				)
			}
		}
	}
}
