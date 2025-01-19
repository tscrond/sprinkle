package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ResourceRepository struct {
	Database *gorm.DB
}

func NewResourceRepository(dbPath string) (*ResourceRepository, error) {

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot initialize without a database")
		return nil, err
	}

	if err := db.AutoMigrate(&HostConfig{}, &MachineConfig{}, &SSHKey{}, &VMConfig{}, &LXCConfig{}); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &ResourceRepository{Database: db}, nil
}
