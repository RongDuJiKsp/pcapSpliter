package flowmaker

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"os"
	"pcapSpliter/utils/packetutil"
)

func MakeSession(file *os.File, dumpPath string) error {
	hd, err := pcap.OpenOfflineFile(file)
	if err != nil {
		return err
	}
	packetSource := gopacket.NewPacketSource(hd, hd.LinkType())
	flows := make(map[string]*flow)
	for packet := range packetSource.Packets() {
		ip, ok := packetutil.AsIpv4(packet)
		if !ok {
			fmt.Println("Find No ipv4 packet")
			continue
		}
		tcp, ok := packetutil.AsTcp(packet)
		if !ok {
			fmt.Println("Find No tcp packet")
			continue
		}
		var key string
		if tcp.SrcPort == 443 {
			key = keyOf(ip, tcp, DirectionReverse)
		} else if tcp.DstPort == 443 {
			key = keyOf(ip, tcp, DirectionForward)
		} else {
			fmt.Println("Find No https packet")
			continue
		}
		flow, ok := flows[key]
		if !ok || flow.syn() && tcp.SYN {
			newerFlow, err := newFlow(packet, dumpPath)
			fmt.Println("New flow", key)
			if err != nil {
				if !errors.Is(err, ErrorNoDoH) {
					fmt.Println(err)
				}
				continue
			}
			if ok {
				flow.close()
			}
			flows[key] = newerFlow
		}
		flows[key].write(packet)
	}
	for _, flow := range flows {
		flow.close()
	}
	return nil
}
