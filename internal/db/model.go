package db

import "gorm.io/gorm"

type HostConfig struct {
	gorm.Model `copier:"-"`
	TargetNode string
	ApiURL     string
	LXCs       []LXCConfig `gorm:"foreignKey:HostID"`
	VMs        []VMConfig  `gorm:"foreignKey:HostID"`
}

type LXCConfig struct {
	gorm.Model `copier:"-"`
	HostID     uint            `copier:"-"`
	Machines   []MachineConfig `gorm:"foreignKey:LXCConfigID"`
}

type VMConfig struct {
	gorm.Model `copier:"-"`
	HostID     uint            `copier:"-"`
	Machines   []MachineConfig `gorm:"foreignKey:LXCConfigID"`
}

type MachineConfig struct {
	gorm.Model       `copier:"-"`
	LXCConfigID      uint `copier:"-"`
	VMConfigID       uint `copier:"-"`
	Name             string
	VmId             int
	OsTemplate       string `gorm:"column:os-template"`
	ISO              string
	IPAddress        string
	CPUs             int
	Memory           int
	DiskSize         int
	SwapSize         int
	Tags             string
	StartOnBoot      bool `gorm:"column:start-on-boot"`
	StorageBackend   string
	TemplateBackend  string
	NetworkBridge    string
	NetworkInterface string
	DefaultGateway   string
	SSHPublicKeys    []SSHKey `gorm:"foreignKey:MachineConfigID"`
}

type SSHKey struct {
	gorm.Model      `copier:"-"`
	MachineConfigID uint `copier:"-"`
	Key             string
	Path            string
}
