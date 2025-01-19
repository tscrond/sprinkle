package config

import "github.com/spf13/viper"

type ConfigManager struct {
	ConfigYAML *HostConfigYAML
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		ConfigYAML: &HostConfigYAML{},
	}
}

func (cm *ConfigManager) LoadConfigFromYAML(configFile string) (*HostConfigYAML, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cm.ConfigYAML); err != nil {
		return nil, err
	}
	config := cm.ConfigYAML

	return config, nil
}
