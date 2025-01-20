package db

import (
	"errors"
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

	if err := db.AutoMigrate(&HostConfig{}, &MachineConfig{}, &SSHKey{}, &VMConfig{}, &LXCConfig{}, &Credentials{}); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &ResourceRepository{Database: db}, nil
}

func (repo *ResourceRepository) SaveCredentials(creds *Credentials) error {
	if err := repo.Database.Create(creds).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ResourceRepository) GetCredentials(targetNode, apiUrl string) (*Credentials, error) {
	var creds Credentials

	if err := repo.Database.Where("target_node = ? AND api_url = ?", targetNode, apiUrl).First(&creds).Error; err != nil {
		return nil, err
	}
	return &creds, nil
}

func (repo *ResourceRepository) CredsExist(targetNode, apiUrl string) bool {
	result := repo.Database.Where("target_node = ? AND api_url = ?", targetNode, apiUrl).First(&Credentials{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		fmt.Println(result.Error)
		return false
	}
	return true
}

func (repo *ResourceRepository) UpdateCredentials(creds *Credentials) error {
	return repo.Database.Where("target_node = ? AND api_url = ?", creds.TargetNode, creds.ApiUrl).First(&Credentials{}).Error
}
