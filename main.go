package main

import (
	"fmt"
	"log"

	"github.com/Soulou/vgdisplay-go/dev"
	"github.com/Soulou/vgdisplay-go/lvm"
)

type LabelHeader struct {
	ID       [8]byte /* LABELONE */
	Sector   uint64  /* Sector number of this label */
	CRC      uint32  /* From next field to end of sector */
	Offset   uint32  /* Offset from start of struct to contents */
	Typename [8]byte /* LVM2 001 */
}

const (
	ID_LEN = 32
)

type explicitPVHeader struct {
	UUID [ID_LEN]byte
	/* This size can be overridden if PV belongs to a VG */
	DeviceSize uint64 /* Bytes */
}

type PVHeader struct {
	explicitPVHeader

	/* NULL-terminated list of data areas followed by */
	/* NULL-terminated list of metadata area headers */
	DiskAreas     []DiskLocn
	MetadataAreas []DiskLocn
}

type explicitPVHeaderExtension struct {
	Version uint32
	Flags   uint32
}

type PVHeaderExtension struct {
	explicitPVHeaderExtension

	/* NULL-terminated list of bootloader areas */
	BooloaderAreas []DiskLocn
}

type DiskLocn struct {
	Offset uint64
	Size   uint64
}

func main() {
	d, err := dev.Open("/dev/loop0")
	if err != nil {
		log.Fatal("fail to open device:", err)
	}
	defer d.Close()

	pv, err := lvm.NewPhysicalVolume(d)
	if err != nil {
		log.Fatal("fail to get PV:", err)
	}
	fmt.Println(pv)
}

// var pvHeader PVHeader
// err = binary.Read(bytes.NewReader(labelSector[header.Offset:]), binary.LittleEndian, &(pvHeader.explicitPVHeader))
// if err != nil {
// 	log.Fatalln("Fail to read PVHeader", err)
// }
// offset := header.Offset + uint32(unsafe.Sizeof(pvHeader.explicitPVHeader))
// areas, shift, err := readDiskLocationList(labelSector[offset:])
// if err != nil {
// 	log.Fatalln("Fail to read disk location list", err)
// }
// pvHeader.DiskAreas = areas
// offset += shift

// areas, shift, err = readDiskLocationList(labelSector[offset:])
// if err != nil {
// 	log.Fatalln("Fail to read disk location list", err)
// }
// pvHeader.MetadataAreas = areas
// offset += shift

// var pvHeaderExt PVHeaderExtension
// err = binary.Read(bytes.NewReader(labelSector[offset:]), binary.LittleEndian, &(pvHeaderExt.explicitPVHeaderExtension))
// if err != nil {
// 	log.Fatalln("Fail to read PVHeaderExtension", err)
// }
// offset = offset + uint32(unsafe.Sizeof(pvHeaderExt.explicitPVHeaderExtension))

// areas, _, err = readDiskLocationList(labelSector[offset:])
// if err != nil {
// 	log.Fatalln("Fail to read disk location list", err)
// }
// pvHeaderExt.BooloaderAreas = areas

// fmt.Printf("PV UUID is %s\n", string(pvHeader.UUID[:]))
// fmt.Println(pvHeader.DiskAreas)
// fmt.Println(pvHeader.MetadataAreas)
// fmt.Printf("%+v\n", pvHeaderExt)

// }

// func readDiskLocationList(buf []byte) ([]DiskLocn, uint32, error) {
// var res []DiskLocn
// var offset uint32 = 0
// for {
// 	var locn DiskLocn
// 	err := binary.Read(bytes.NewReader(buf[offset:]), binary.LittleEndian, &locn)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	offset += uint32(unsafe.Sizeof(locn))
// 	if locn.Offset == 0 {
// 		break
// 	}
// 	res = append(res, locn)
// }
// return res, offset, nil
// }
