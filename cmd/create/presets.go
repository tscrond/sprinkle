package create

import (
	"github.com/tscrond/sprinkle/lib"
)

var SMALL_PRESET_WORKER = lib.MachineConfig{
	Name:     "worker-node-small",
	DiskSize: 30,
	SwapSize: 0,
	CPUCount: 2,
	Memory:   2048,
	OnBoot:   false,
	Tags:     "k8s;worker;small",
}

var MEDIUM_PRESET_WORKER = lib.MachineConfig{
	Name:     "worker-node-medium",
	DiskSize: 60,
	SwapSize: 0,
	CPUCount: 4,
	Memory:   4096,
	OnBoot:   false,
	Tags:     "k8s;worker;medium",
}

var MEDIUM_STORAGE_PRESET_WORKER = lib.MachineConfig{
	Name:     "worker-node-medium",
	DiskSize: 200,
	SwapSize: 0,
	CPUCount: 4,
	Memory:   4096,
	OnBoot:   false,
	Tags:     "k8s;worker;medium;storage",
}

var LARGE_PRESET_WORKER = lib.MachineConfig{
	Name:     "worker-node-large",
	DiskSize: 120,
	SwapSize: 0,
	CPUCount: 8,
	Memory:   8192,
	OnBoot:   false,
	Tags:     "k8s;worker;large",
}

var SMALL_PRESET_MASTER = lib.MachineConfig{
	Name:     "master-node-small",
	DiskSize: 30,
	SwapSize: 0,
	CPUCount: 2,
	Memory:   2048,
	OnBoot:   false,
	Tags:     "k8s;master;small",
}

var SMALL_STORAGE_PRESET_MASTER = lib.MachineConfig{
	Name:     "master-node-small",
	DiskSize: 200,
	SwapSize: 0,
	CPUCount: 2,
	Memory:   2048,
	OnBoot:   false,
	Tags:     "k8s;master;small;storage",
}

var MEDIUM_PRESET_MASTER = lib.MachineConfig{
	Name:     "master-node-medium",
	DiskSize: 60,
	SwapSize: 0,
	CPUCount: 4,
	Memory:   4096,
	OnBoot:   false,
	Tags:     "k8s;master;medium",
}

var LARGE_PRESET_MASTER = lib.MachineConfig{
	Name:     "master-node-large",
	DiskSize: 120,
	SwapSize: 0,
	CPUCount: 8,
	Memory:   8192,
	OnBoot:   false,
	Tags:     "k8s;worker;large",
}
