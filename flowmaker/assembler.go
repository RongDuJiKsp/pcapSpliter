package flowmaker

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"pcapSpliter/utils/dns"
	"pcapSpliter/utils/packetutil"
	"time"
)

var (
	ErrorNoIpv4  = errors.New("invalid ipv4 packet")
	ErrorNoIpv6  = errors.New("invalid ipv6 packet")
	ErrorNoTcp   = errors.New("invalid tcp packet")
	ErrorNoUdp   = errors.New("invalid udp packet")
	ErrorNoDoH   = errors.New("invalid doh packet")
	ErrorNoHttps = errors.New("invalid https packet")
)

const (
	DirectionForward = "forward"
	DirectionReverse = "reverse"
)
const TimeFmtForFile = "2006-01-02_Mon_15-04-05_MST"

func capInfo(baseName string, first gopacket.Packet) (fileName string, firstFlow time.Time, providerIP string, key string, err error) {
	firstFlow = first.Metadata().CaptureInfo.Timestamp
	ip, ok := packetutil.AsIpv4(first)
	if !ok {
		return "", firstFlow, "", "", ErrorNoIpv4
	}
	tcp, ok := packetutil.AsTcp(first)
	if !ok {
		return "", firstFlow, "", "", ErrorNoTcp
	}
	if tcp.SrcPort == 443 {
		providerIP = ip.SrcIP.String()
		key = keyOf(ip, tcp, DirectionReverse)
	} else if tcp.DstPort == 443 {
		providerIP = ip.DstIP.String()
		key = keyOf(ip, tcp, DirectionForward)
	} else {
		return "", firstFlow, "", "", ErrorNoHttps
	}
	provider := dns.ConvertIPToDNSProvider(providerIP)
	if provider == dns.ProviderUnknown {
		return "", firstFlow, "", "", ErrorNoDoH
	}
	fileName = fmt.Sprintf("%s_%s_%s.pcap", baseName, provider, firstFlow.Format(TimeFmtForFile))
	return
}

func keyOf(ip *layers.IPv4, tcp *layers.TCP, direction string) string {
	var src, dst string
	if direction == DirectionForward {
		src = fmt.Sprintf("%s:%d", ip.SrcIP.String(), tcp.SrcPort)
		dst = fmt.Sprintf("%s:%d", ip.DstIP.String(), tcp.DstPort)
	} else if direction == DirectionReverse {
		src = fmt.Sprintf("%s:%d", ip.DstIP.String(), tcp.DstPort)
		dst = fmt.Sprintf("%s:%d", ip.SrcIP.String(), tcp.SrcPort)
	}
	return fmt.Sprintf("%s->%s", src, dst)
}
