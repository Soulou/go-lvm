package dev

import "os"

type Device struct {
	*os.File
	BlockSize         uint32
	PhysicalBlockSize uint32
}
