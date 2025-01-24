package db

import "time"

type HostConfig struct {
	Metadata
	TargetNode string `gorm:"primaryKey"`
	ApiURL     string
	Machines   []MachineConfig `gorm:"foreignKey:TargetNode;constraint:OnDelete:CASCADE"`
}

type MachineConfig struct {
	Metadata
	VmId             int `gorm:"primaryKey"` // Primary Key
	Type             string
	Name             string
	OsTemplate       string
	ISO              string
	IPAddress        string
	CPUs             int
	Memory           int
	DiskSize         int
	SwapSize         int
	Tags             string
	StartOnBoot      bool
	StorageBackend   string
	TemplateBackend  string
	NetworkBridge    string
	NetworkInterface string
	DefaultGateway   string
	TargetNode       string   `gorm:"index"` // Foreign Key to HostConfig
	SSHPublicKeys    []SSHKey `gorm:"foreignKey:VmId;constraint:OnDelete:CASCADE"`
}

type SSHKey struct {
	Metadata
	ID   uint    `gorm:"primaryKey"`
	VmId int     `gorm:"index"`
	Key  *string `gorm:"type:text"`
	Path *string `gorm:"type:text"`
}

type Credentials struct {
	Metadata
	TargetNode string
	ApiUrl     string
	Username   string
	Password   string
	CsrfToken  string
	PVETicket  string
}

type Metadata struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
