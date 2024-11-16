package packetutil

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func AsIpv4(p gopacket.Packet) (*layers.IPv4, bool) {
	ipv4LayerPoint := p.Layer(layers.LayerTypeIPv4)
	if ipv4LayerPoint == nil {
		return nil, false
	}
	ipv4Layer := ipv4LayerPoint.(*layers.IPv4)
	return ipv4Layer, true
}
func AsTcp(p gopacket.Packet) (*layers.TCP, bool) {
	tcpLayerPoint := p.Layer(layers.LayerTypeTCP)
	if tcpLayerPoint == nil {
		return nil, false
	}
	tcpLayer := tcpLayerPoint.(*layers.TCP)
	return tcpLayer, true
}

func AsUdp(p gopacket.Packet) (*layers.UDP, bool) {
	udpLayerPoint := p.Layer(layers.LayerTypeUDP)
	if udpLayerPoint == nil {
		return nil, false
	}
	udpLayer := udpLayerPoint.(*layers.UDP)
	return udpLayer, true
}
