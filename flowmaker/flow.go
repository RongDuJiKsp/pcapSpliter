package flowmaker

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
	"os"
	"path"
	"path/filepath"
	"time"
)

const FlowMinSize = 8

type flow struct {
	packetWriter *pcapgo.Writer
	pcapFile     *os.File
	firstFlow    time.Time
	providerIP   string
	key          string
	size         uint
}

func newFlow(firstPkg gopacket.Packet, dumpPath string) (*flow, error) {
	fileName, firstFlow, providerIP, key, err := capInfo(filepath.Base(dumpPath), firstPkg)
	if err != nil {
		return nil, err
	}
	fmt.Println("Open:", fileName)
	file, err := os.Create(path.Join(dumpPath, fileName))
	if err != nil {
		return nil, err
	}
	w := pcapgo.NewWriter(file)
	if err = w.WriteFileHeader(1500, layers.LinkTypeEthernet); err != nil {
		_ = file.Close()
		return nil, err
	}
	return &flow{packetWriter: w, pcapFile: file, firstFlow: firstFlow, providerIP: providerIP, key: key}, err

}

func (f *flow) write(p gopacket.Packet) {
	_ = f.packetWriter.WritePacket(p.Metadata().CaptureInfo, p.Data())
	f.size++
}
func (f *flow) syn() bool {
	return f.size > FlowMinSize
}

func (f *flow) close() {
	f.pcapFile.Close()
}
