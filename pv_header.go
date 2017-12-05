package lvm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/Soulou/go-lvm/crc32"

	"github.com/pkg/errors"
)

const (
	LabelID = "LABELONE"
	// Label can be in any of the first 4 sectors
	LabelScanSectors = 4
	LabelScanSize    = LabelScanSectors * SectorSize
	LabelHeaderSize  = 32

	PhysicalVolumeIDLength = 32
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

func (h LabelHeader) CheckCRC32(sector []byte) (uint32, bool) {
	crc32checksum := crc32.Calc(InitialCRC, sector[20:SectorSize])
	return crc32checksum, crc32checksum == h.CRC
}

type StaticPhysicalVolumeHeader struct {
	UUID [PhysicalVolumeIDLength]byte
	/* This size can be overridden if PV belongs to a VG */
	DeviceSize uint64 /* Bytes */
}

type PhysicalVolumeHeader struct {
	StaticPhysicalVolumeHeader
	/* NULL-terminated list of data areas followed by */
	/* NULL-terminated list of metadata area headers */
	DiskAreas     []DataAreaDescriptor
	MetadataAreas []DataAreaDescriptor
}

func (h PhysicalVolumeHeader) UUIDToString() string {
	buf := bytes.NewBuffer(make([]byte, 0, 38))
	buf.Write(h.UUID[0:6])
	buf.Write([]byte{'-'})
	buf.Write(h.UUID[6:10])
	buf.Write([]byte{'-'})
	buf.Write(h.UUID[14:18])
	buf.Write([]byte{'-'})
	buf.Write(h.UUID[18:22])
	buf.Write([]byte{'-'})
	buf.Write(h.UUID[22:26])
	buf.Write([]byte{'-'})
	buf.Write(h.UUID[26:32])
	return buf.String()
}

func (h PhysicalVolumeHeader) String() string {
	return fmt.Sprintf(
		"PVHeader[ID:%s Size:%d DiskAreas:%d MetadataAreas:%d]",
		h.UUIDToString(), h.DeviceSize, len(h.DiskAreas), len(h.MetadataAreas),
	)
}

type StaticPhysicalVolumeHeaderExtension struct {
	Version uint32
	Flags   uint32
}

type PhysicalVolumeHeaderExtension struct {
	StaticPhysicalVolumeHeaderExtension

	/* NULL-terminated list of bootloader areas */
	BooloaderAreas []DataAreaDescriptor
}

func (ext PhysicalVolumeHeaderExtension) String() string {
	return fmt.Sprintf("Ext[Version:%d Flags:%d, BootloaderAreas:%d]", ext.Version, ext.Flags, len(ext.BooloaderAreas))
}

func (pv *PhysicalVolume) readHeaderExt(reader io.Reader) (PhysicalVolumeHeaderExtension, error) {
	var ext PhysicalVolumeHeaderExtension
	err := binary.Read(reader, binary.LittleEndian, &(ext.StaticPhysicalVolumeHeaderExtension))
	if err != nil {
		return ext, errors.Wrapf(err, "fail to parse pv header extension")
	}

	areas, err := readDataAreaDescriptorList(reader)
	if err != nil {
		return ext, errors.Wrapf(err, "fail to parse data area descriptors for bootloader areas")
	}
	ext.BooloaderAreas = areas

	return ext, nil
}

func (pv *PhysicalVolume) readHeader(reader io.Reader) (PhysicalVolumeHeader, error) {
	var header PhysicalVolumeHeader
	err := binary.Read(reader, binary.LittleEndian, &(header.StaticPhysicalVolumeHeader))
	if err != nil {
		return header, errors.Wrapf(err, "fail to parse pv header")
	}

	areas, err := readDataAreaDescriptorList(reader)
	if err != nil {
		return header, errors.Wrapf(err, "fail to parse data area descriptors for disk areas")
	}
	header.DiskAreas = areas

	areas, err = readDataAreaDescriptorList(reader)
	if err != nil {
		return header, errors.Wrapf(err, "fail to parse data area descriptors for metadata areas")
	}
	header.MetadataAreas = areas

	return header, nil
}

// It is expected to get the first block (block n°0) of the device
func (pv *PhysicalVolume) readLabelHeader(reader io.Reader) (LabelHeader, error) {
	var header LabelHeader

	for sector := 0; sector < LabelScanSectors; sector++ {
		err := binary.Read(reader, binary.LittleEndian, &header)
		if err != nil {
			log.Println("error reading error: ", err, "continuing..")
			continue
		}

		if string(header.ID[:]) == LabelID {
			if header.Sector != uint64(sector) {
				return header, errors.Errorf("header sector does not match: (%v, expected: %v)", header.Sector, sector)
			}
			break
		} else {
			toNextSector := make([]byte, SectorSize-LabelHeaderSize)
			_, err := reader.Read(toNextSector)
			if err != nil {
				return header, errors.Wrapf(err, "fail to read to next sector")
			}
		}
	}
	if header.Sector == 0 {
		return header, errors.New("label header not found")
	}
	return header, nil
}
