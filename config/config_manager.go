package config

import (
	"github.com/spf13/viper"
)

type ConfigManager struct {
	ConfigYAML *HostConfigYAML
	DbPath     string
}

func NewConfigManager(dbPath string) *ConfigManager {
	return &ConfigManager{
		ConfigYAML: &HostConfigYAML{},
		DbPath:     dbPath,
	}
}

func (cm *ConfigManager) LoadConfigFromYAML() (*HostConfigYAML, error) {
	viper.SetConfigFile(cm.DbPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cm.ConfigYAML); err != nil {
		return nil, err
	}
	config := cm.ConfigYAML

	return config, nil
}
