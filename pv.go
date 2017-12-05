package lvm

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/Soulou/go-lvm/dev"
	"github.com/pkg/errors"
)

type PhysicalVolume struct {
	device          *dev.Device
	labelHeader     LabelHeader
	header          PhysicalVolumeHeader
	headerExt       PhysicalVolumeHeaderExtension
	metadataHeaders []MetadataHeader
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

	reader := bytes.NewReader(b)
	labelHeader, err := pv.readLabelHeader(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read label header")
	}
	if crc32checksum, ok := labelHeader.CheckCRC32(b[labelHeader.Sector*SectorSize:]); !ok {
		log.Printf("Fail to check label header checksum, got %v, expected: %v", crc32checksum, labelHeader.CRC)
	}
	pv.labelHeader = labelHeader

	pvHeader, err := pv.readHeader(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read pv header")
	}
	pv.header = pvHeader

	pvHeaderExt, err := pv.readHeaderExt(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to reade pv header ext")
	}
	pv.headerExt = pvHeaderExt

	metadataHeaders, err := pv.readMetadataHeaders()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read metadata header areas")
	}
	pv.metadataHeaders = metadataHeaders

	fmt.Println(labelHeader)
	fmt.Println(pvHeader)
	fmt.Println(pvHeaderExt)
	fmt.Println(metadataHeaders)

	for _, h := range metadataHeaders {
		m, err := pv.readMetadata(h)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to read pv metadata")
		}
		fmt.Println(string(m))
	}

	return pv, nil
}

func (pv *PhysicalVolume) ReadBlock(offset uint64) ([]byte, error) {
	block := make([]byte, pv.device.BlockSize)
	_, err := pv.device.Seek(int64(offset), os.SEEK_SET)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to seek in device")
	}
	_, err = pv.device.Read(block)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read from device")
	}
	return block, nil
}
