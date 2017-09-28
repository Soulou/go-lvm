package lvm

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

type DataAreaDescriptor struct {
	Offset uint64
	Size   uint64
}

func readDataAreaDescriptorList(reader io.Reader) ([]DataAreaDescriptor, error) {
	var res []DataAreaDescriptor
	for {
		var desc DataAreaDescriptor
		err := binary.Read(reader, binary.LittleEndian, &desc)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to parse data area descriptor")
		}
		if desc.Offset == 0 {
			break
		}
		res = append(res, desc)
	}
	return res, nil
}
