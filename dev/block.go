package dev

import (
	"syscall"
	"unsafe"
)

func (d *Device) ReadBlockSize() (uint32, error) {
	var blksize uint32
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.File.Fd(), uintptr(BLKBSZGET), uintptr(unsafe.Pointer(&blksize)))
	if errno != 0 {
		return 0, errno
	}

	return blksize, nil
}

func (d *Device) ReadPhysicalBlockSize() (uint32, error) {
	var physicalBlksize uint32
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, d.File.Fd(), uintptr(BLKPBSZGET), uintptr(unsafe.Pointer(&physicalBlksize)))
	if errno != 0 {
		return 0, errno
	}
	return physicalBlksize, nil
}
