package config

import "github.com/spf13/viper"

type ConfigManager struct {
	Config *HostConfigYAML
}

func NewConfigManager() *ConfigManager {

	return &ConfigManager{
		Config: &HostConfigYAML{},
	}
}

func (cm *ConfigManager) LoadConfigFromYAML(configFile string) error {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(cm.Config); err != nil {
		return err
	}

	return nil
}
