package db

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (repo *ResourceRepository) GetAllHostConfigs() ([]HostConfig, error) {
	var hostConfigs []HostConfig

	// Use Preload to fetch all related nested data
	if err := repo.Database.
		Preload("LXCs.Machines.SSHPublicKeys").
		Preload("VMs.Machines.SSHPublicKeys").
		Find(&hostConfigs).Error; err != nil {
		return nil, err
	}

	return hostConfigs, nil
}

func (repo *ResourceRepository) InsertHostConfigs(hostConfigs []HostConfig) error {
	return repo.Database.Transaction(func(tx *gorm.DB) error {
		for _, hostConfig := range hostConfigs {
			// Insert or update each HostConfig
			if err := tx.Clauses(clause.OnConflict{
				UpdateAll: true, // Update all fields if a conflict occurs
			}).Create(&hostConfig).Error; err != nil {
				return fmt.Errorf("failed to insert/update host config for node %s: %w", hostConfig.TargetNode, err)
			}

			// Insert child LXCs for this HostConfig
			for _, lxc := range hostConfig.LXCs {
				lxc.HostID = hostConfig.ID
				if err := tx.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&lxc).Error; err != nil {
					return fmt.Errorf("failed to insert/update LXC config for host %s: %w", hostConfig.TargetNode, err)
				}

				// Insert Machines under this LXC
				for _, machine := range lxc.Machines {
					machine.LXCConfigID = lxc.ID
					if err := tx.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&machine).Error; err != nil {
						return fmt.Errorf("failed to insert/update machine under LXC %s: %w", machine.Name, err)
					}
				}
			}

			// Insert child VMs for this HostConfig
			for _, vm := range hostConfig.VMs {
				vm.HostID = hostConfig.ID
				if err := tx.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&vm).Error; err != nil {
					return fmt.Errorf("failed to insert/update VM config for host %s: %w", hostConfig.TargetNode, err)
				}

				// Insert Machines under this VM
				for _, machine := range vm.Machines {
					machine.VMConfigID = vm.ID
					if err := tx.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&machine).Error; err != nil {
						return fmt.Errorf("failed to insert/update machine under VM %s: %w", machine.Name, err)
					}
				}
			}
		}

		return nil
	})
}
