package lvm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"

	"github.com/pkg/errors"
)

const (
	MetadataHeaderSize = SectorSize
)

type StaticMetadataHeader struct {
	// Checksum from offset 4 to end of the header
	CRC uint32
	// "\x20LVM2\x20x[5A%r0N*>"
	Signature [16]byte
	// 1
	Version uint32
	// Offset from the beginning of the disk to the metadata header area
	Offset uint64
	// Size of the metadata area
	Size uint64
}

type MetadataHeader struct {
	StaticMetadataHeader
	Locations []RawLocationDescriptor
}

func (h MetadataHeader) String() string {
	return fmt.Sprintf(
		"MetadataHeader[CRC:%d Signature:%s Version:%d Offset:%d Size:%d Locations:%d]",
		h.CRC, string(h.Signature[:]), h.Version, h.Offset, h.Size, len(h.Locations),
	)
}

func (h MetadataHeader) CheckCRC32(block []byte) (uint32, bool) {
	headerBytes := block[4:MetadataHeaderSize]
	crc32checksum := crc32.Update(InitialCRC, crcTable, headerBytes)
	return crc32checksum, crc32checksum == h.CRC
}

func (pv *PhysicalVolume) readMetadataHeaders() ([]MetadataHeader, error) {
	var headers []MetadataHeader
	for _, area := range pv.header.MetadataAreas {
		b, err := pv.ReadBlock(area.Offset)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to read device block for metadata area")
		}
		reader := bytes.NewReader(b)

		header, err := pv.readMetadataHeader(reader)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to read metadata header")
		}

		crc32checksum, ok := header.CheckCRC32(b)
		if !ok {
			log.Println("Fail to check metadata header checksum, got", crc32checksum, "expected", header.CRC)
		}

		headers = append(headers, header)
	}
	return headers, nil
}

func (pv *PhysicalVolume) readMetadataHeader(reader io.Reader) (MetadataHeader, error) {
	var header MetadataHeader
	err := binary.Read(reader, binary.LittleEndian, &(header.StaticMetadataHeader))
	if err != nil {
		return header, errors.Wrapf(err, "fail to parse pv header")
	}

	locations, err := readRawLocationDescriptorList(reader)
	if err != nil {
		return header, errors.Wrapf(err, "fail to parse data area descriptors for disk areas")
	}
	header.Locations = locations
	return header, nil
}

func (pv *PhysicalVolume) readMetadata(h MetadataHeader) ([]byte, error) {
	for _, loc := range h.Locations {
		block, err := pv.ReadBlock(h.Offset + loc.Offset)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to read block")
		}
		return block, nil
	}
	return nil, nil
}
