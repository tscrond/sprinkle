package config

import (
	"strconv"

	"golang.org/x/exp/rand"
)

// random number from 100 to 999999999
var randomNumber = rand.Intn(999999999-100+1) + 100
var randomNumberStr = strconv.Itoa(randomNumber)

var DEFAULT_MACHINE_CONFIG = MachineConfigYAML{
	Name:             "machine-" + randomNumberStr,
	VmId:             randomNumber,
	OsTemplate:       "debian-11-standard_11.7-1_amd64.tar.zst",
	NetworkBridge:    "vmbr0",
	NetworkInterface: "eth0",
	DefaultGateway:   "192.168.1.1",
	IPAddress:        "192.168.1.150/24",
	StorageBackend:   "local-lvm",
	TemplateBackend:  "local",
	DiskSize:         30,
	SwapSize:         0,
	CPUs:             2,
	StartOnBoot:      false,
	ISO:              "ubuntu-22.04.3-live-server-amd64.iso",
	Tags:             "",
}
