package lvm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"

	"github.com/juju/errgo/errors"
)

const (
	LabelID = "LABELONE"
	// Label can be in any of the first 4 sectors
	LabelScanSectors = 4
	LabelScanSize    = LabelScanSectors * SectorSize
)

type LabelHeader struct {
	ID       [8]byte /* LABELONE */
	Sector   uint64  /* Sector number of this label */
	CRC      uint32  /* From next field to end of sector */
	Offset   uint32  /* Offset from start of struct to contents */
	Typename [8]byte /* LVM2 001 */
}

func (h LabelHeader) String() string {
	return fmt.Sprintf(
		"LabelHeader[ID:%s Sector:%d CRC:%d Offset:%d Typename:%s]",
		string(h.ID[:]), h.Sector, h.CRC, h.Offset, string(h.Typename[:]),
	)
}

// It is expected to get the first block (block n°0) of the device
func (pv *PhysicalVolume) readLabelHeader(block []byte) (*LabelHeader, error) {
	var header LabelHeader

	for sector := 0; sector < LabelScanSectors; sector++ {
		offset := sector * SectorSize
		buffer := block[offset : offset+SectorSize]
		err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &header)
		if err != nil {
			log.Println("error reading error: ", err, "continuing..")
			continue
		}

		if string(header.ID[:]) == LabelID {
			if header.Sector != uint64(sector) {
				return nil, errors.Newf("header sector does not match: (%v, expected: %v)", header.Sector, sector)
			}

			crc32q := crc32.MakeTable(InitialCRC)
			crc32checksum := crc32.Checksum(buffer[20:], crc32q)
			if crc32checksum != header.CRC {
				log.Printf("Checksum doesn't match :( (%v, expected: %v)", crc32checksum, header.CRC)
			}
			break
		}
	}
	if header.Sector == 0 {
		return nil, errors.New("label header not found")
	}
	return &header, nil
}
