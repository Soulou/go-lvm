package dev

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

func Open(path string) (*Device, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY|syscall.O_DIRECT|syscall.O_NOATIME, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to open %s", path)
	}

	d := &Device{File: fd}

	blksize, err := d.ReadBlockSize()
	if err != nil {
		d.Close()
		return nil, errors.Wrapf(err, "fail to read block size")
	}
	d.BlockSize = blksize

	pblksize, err := d.ReadPhysicalBlockSize()
	if err != nil {
		d.Close()
		return nil, errors.Wrapf(err, "fail to read physical block size")
	}
	d.PhysicalBlockSize = pblksize

	return d, nil
}
