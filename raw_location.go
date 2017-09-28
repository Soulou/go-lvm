package lvm

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	RawLocationDescriptorSize = 24
)

type RawLocationDescriptor struct {
	// Offset to the beginning of the Location
	Offset uint64
	// Size of the location
	Size uint64
	// CRC sum of ???
	CRC uint32
	// 1: Location should be ignored
	Flags uint32
}

func readRawLocationDescriptorList(reader io.Reader) ([]RawLocationDescriptor, error) {
	var res []RawLocationDescriptor
	for {
		var loc RawLocationDescriptor
		err := binary.Read(reader, binary.LittleEndian, &loc)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to parse raw location descriptor")
		}
		if loc.Offset == 0 {
			break
		}
		res = append(res, loc)
	}
	return res, nil
}
