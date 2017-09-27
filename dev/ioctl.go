package dev

import "unsafe"

// From asm/ioctl.h
const (
	_IOC_NRBITS   = 8
	_IOC_TYPEBITS = 8
	_IOC_SIZEBITS = 14
	_IOC_DIRBITS  = 2

	_IOC_NRSHIFT   = 0
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS

	_IOC_NONE = 0
	_IOC_READ = 2
)

func _IO(t, nr int) int {
	return _IOC(_IOC_NONE, t, nr, 0)
}

func _IOC(dir, t, nr int, size uintptr) int {
	return (dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) | (nr << _IOC_NRSHIFT) | (int(size) << _IOC_SIZESHIFT)
}

func _IOR(t, nr int, size uintptr) int {
	return _IOC(_IOC_READ, t, nr, unsafe.Sizeof(size))
}

var (
	BLKBSZGET  = _IOR(0x12, 112, uintptr(unsafe.Pointer((*byte)(nil))))
	BLKPBSZGET = _IO(0x12, 123)
)
