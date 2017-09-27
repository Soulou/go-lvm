package lvm

import (
	"fmt"
	"os"

	"github.com/Soulou/vgdisplay-go/dev"
	"github.com/pkg/errors"
)

type PhysicalVolume struct {
	device *dev.Device
	header *LabelHeader
}

func NewPhysicalVolume(dev *dev.Device) (*PhysicalVolume, error) {
	if dev == nil {
		return nil, errors.New("nil device")
	}
	pv := &PhysicalVolume{device: dev}

	b, err := pv.ReadBlock(0)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read block 0")
	}

	header, err := pv.readLabelHeader(b)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read label header")
	}
	pv.header = header
	fmt.Println(header)

	return pv, nil
}

func (pv *PhysicalVolume) ReadBlock(n uint32) ([]byte, error) {
	block := make([]byte, pv.device.BlockSize)
	_, err := pv.device.Seek(int64(n*pv.device.BlockSize), os.SEEK_SET)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to seek in device")
	}
	_, err = pv.device.Read(block)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read from device")
	}
	return block, nil
}
