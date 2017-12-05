package metadata

// pv0 {
//   id = "Foeuwm-r0Pm-dbMn-dn3a-Gkfc-ZQvh-naxBQ6"
//   device = "/dev/loop0"
//   status = ["ALLOCATABLE"]
//   flags = []
//   dev_size = 20971520
//   pe_start = 2048
//   pe_count = 2559
// }
type PhysicalVolume struct {
	ID      string
	Device  string
	Status  []string
	Flags   []string
	DevSize int
	PeStart int
	PeCount int
}

// stripes = [
//   "pv0", 0
// ]
type Stripe struct {
	Name   string
	Offset int64
}

// segment1 {
//   start_extent = 0
//   extent_count = 1

//   type = "striped"
//   stripe_count = 1

//   stripes = [...]
// }
type Segment struct {
	// Commong in LVs
	StartExtent   int
	ExtentCount   int
	TransactionID int
	Type          string

	// In a thin pool
	// Type == "thin-pool"
	Metadata      string
	Pool          string
	ChunkSize     int
	Discards      string
	ZeroNewBlocks int

	// In a thin LV
	// Type == "thin"
	ThinPool string
	DeviceID int

	// In a thin-pool component
	// Type == "striped"
	StripeCount int
	Stripes     []Stripe
}

// testlv {
//   id = "P9O9Oi-veHY-ENNF-eyoZ-gabQ-YvbQ-NUzGbw"
//   status = ["READ", "WRITE", "VISIBLE"]
//   flags = []
//   creation_time = 1507908271
//   creation_host = "soulou7"
//   segment_count = 1
//   segment1 { ... }
// }
type LogicalVolume struct {
	Name         string
	ID           string
	Status       []string
	Flags        []string
	CreationTime int64
	CreationHost string
	SegmentCount int

	Segments []Segment
}

// test1 {
//   id = "McoYUJ-C0Pm-FoF9-1maj-tivD-X1ha-7m9Xbc"
//   seqno = 2
//   format = "lvm2"
//   status = ["RESIZEABLE", "READ", "WRITE"]
//   flags = []
//   extent_size = 8192
//   max_lv = 0
//   max_pv = 0
//   metadata_copies = 0
//   physical_volumes { ... }
//   logical_volumes { ... }
// }
//
// contents = "Text Format Volume Group"
// version = 1
// description = ""
// creation_host = "soulou7"       # Linux soulou7 4.13.4-1-ARCH #1 SMP PREEMPT Thu Sep 28 08:39:52 CEST 2017 x86_64
// creation_time = 1507908271      # Fri Oct 13 17:24:31 2017
//}
type VolumeGroup struct {
	Name           string
	ID             string
	Seqno          int
	Format         string
	Status         []string
	Flags          []string
	ExtentSize     int
	MaxLV          int
	MaxPV          int
	MetadataCopies int

	PhysicalVolumes []PhysicalVolume
	LogicalVolumes  []LogicalVolume

	Contents     string
	Version      int
	Description  string
	CreationHost string
	CreationTime int64
}
