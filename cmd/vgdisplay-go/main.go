package main

import (
	"log"

	"github.com/Soulou/go-lvm"
	"github.com/Soulou/go-lvm/dev"
)

func main() {
	d, err := dev.Open("/dev/loop0")
	if err != nil {
		log.Fatal("fail to open device: ", err)
	}
	defer d.Close()

	_, err = lvm.NewPhysicalVolume(d)
	if err != nil {
		log.Fatal("fail to get PV: ", err)
	}

	d, err = dev.Open("/dev/loop1")
	if err != nil {
		log.Fatal("fail to open device: ", err)
	}
	defer d.Close()

	_, err = lvm.NewPhysicalVolume(d)
	if err != nil {
		log.Fatal("fail to get PV: ", err)
	}
}
