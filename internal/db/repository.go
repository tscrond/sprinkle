package db

import (
	"errors"
	"fmt"
	"reflect"
	"time"

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

	if err := db.AutoMigrate(&HostConfig{}, &MachineConfig{}, &SSHKey{}, &Credentials{}); err != nil {
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

func (repo *ResourceRepository) InsertHostConfigs(hostConfigs []HostConfig) error {

	return repo.Database.Transaction(func(tx *gorm.DB) error {
		for _, hostConfig := range hostConfigs {
			// TODO insert host config to DB
			// fmt.Printf("%+v\n", hostConfig)
			if err := repo.InsertOrModifyStruct(&hostConfig, &HostConfig{TargetNode: hostConfig.TargetNode}); err != nil {
				return err
			}
			for _, machine := range hostConfig.Machines {
				if err := repo.InsertOrModifyStruct(&machine, &MachineConfig{VmId: machine.VmId}); err != nil {
					return err
				}
				for _, sshkey := range machine.SSHPublicKeys {
					sshkey.ID = 0
					sshkey.CreatedAt = time.Time{}
					sshkey.UpdatedAt = time.Time{}

					if err := repo.InsertOrModifyStruct(&sshkey, &SSHKey{VmId: machine.VmId, Key: sshkey.Key, Path: sshkey.Path}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

// obj: database model which is being operated on, condition: struct with specified keys on which operation is based
func (repo *ResourceRepository) InsertOrModifyStruct(obj interface{}, condition interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("obj must be a pointer to a struct")
	}

	// fmt.Printf("%+v\n", obj)
	err := repo.Database.Where(condition).FirstOrCreate(obj).Error
	if err != nil {
		return fmt.Errorf("error occurred: %s", err)
	}

	// If the record already exists, update it
	if err := repo.Database.Model(obj).Updates(obj).Error; err != nil {
		return fmt.Errorf("error occurred while updating: %s", err)
	}

	return nil
}

func (repo *ResourceRepository) GetAllHostConfigs() ([]HostConfig, error) {
	var hostConfigs []HostConfig

	// Use Preload to fetch all related nested data
	if err := repo.Database.
		Preload("Machines.SSHPublicKeys").
		Find(&hostConfigs).Error; err != nil {
		return nil, err
	}

	return hostConfigs, nil
}

func (repo *ResourceRepository) CheckIfRecordExists(model interface{}, condition map[string]interface{}) (bool, error) {
	err := repo.Database.Where(condition).First(model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}
