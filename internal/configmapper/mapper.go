package configmapper

import (
	"github.com/tscrond/sprinkle/config"
	"github.com/tscrond/sprinkle/internal/auth"
	"github.com/tscrond/sprinkle/internal/db"
)

func MapConfigToAuthConfig(cfg *config.HostConfigYAML) []auth.AuthConfig {
	var authConfigs []auth.AuthConfig
	for _, host := range cfg.Hosts {
		authConfigs = append(authConfigs, auth.AuthConfig{ApiUrl: host.ApiUrl, TargetNode: host.TargetNode})
	}

	return authConfigs
}

func ConvertConfigToDBModel(cfg *config.HostConfigYAML) ([]db.HostConfig, error) {
	dbHostConfigs := []db.HostConfig{}

	for nodeName, host := range cfg.Hosts {
		hostconfig := db.HostConfig{
			TargetNode: nodeName,
			ApiURL:     host.ApiUrl,
		}

		for _, machine := range host.LXCs.Machines {
			machinePubKeys := []db.SSHKey{}
			for _, pubkey := range machine.SshPublicKeys {
				machinePubKeys = append(machinePubKeys, db.SSHKey{
					VmId: machine.VmId,
					Key:  &pubkey.Key,
					Path: &pubkey.Path,
				})
			}

			hostconfig.Machines = append(hostconfig.Machines, db.MachineConfig{
				Type:             "lxc",
				Name:             machine.Name,
				VmId:             machine.VmId,
				OsTemplate:       machine.OsTemplate,
				ISO:              machine.ISO,
				IPAddress:        machine.IPAddress,
				CPUs:             machine.CPUs,
				Memory:           machine.Memory,
				DiskSize:         machine.DiskSize,
				SwapSize:         machine.SwapSize,
				Tags:             machine.Tags,
				StartOnBoot:      machine.StartOnBoot,
				StorageBackend:   machine.StorageBackend,
				TemplateBackend:  machine.TemplateBackend,
				NetworkBridge:    machine.NetworkBridge,
				NetworkInterface: machine.NetworkInterface,
				DefaultGateway:   machine.DefaultGateway,
				SSHPublicKeys:    machinePubKeys,
			})
			dbHostConfigs = append(dbHostConfigs, hostconfig)
		}

		for _, machine := range host.VMs.Machines {
			machinePubKeys := []db.SSHKey{}
			for _, pubkey := range machine.SshPublicKeys {
				machinePubKeys = append(machinePubKeys, db.SSHKey{
					VmId: machine.VmId,
					Key:  &pubkey.Key,
					Path: &pubkey.Path,
				})
			}

			hostconfig.Machines = append(hostconfig.Machines, db.MachineConfig{
				Type:             "vm",
				Name:             machine.Name,
				VmId:             machine.VmId,
				OsTemplate:       machine.OsTemplate,
				ISO:              machine.ISO,
				IPAddress:        machine.IPAddress,
				CPUs:             machine.CPUs,
				Memory:           machine.Memory,
				DiskSize:         machine.DiskSize,
				SwapSize:         machine.SwapSize,
				Tags:             machine.Tags,
				StartOnBoot:      machine.StartOnBoot,
				StorageBackend:   machine.StorageBackend,
				TemplateBackend:  machine.TemplateBackend,
				NetworkBridge:    machine.NetworkBridge,
				NetworkInterface: machine.NetworkInterface,
				DefaultGateway:   machine.DefaultGateway,
				SSHPublicKeys:    machinePubKeys,
			})
			dbHostConfigs = append(dbHostConfigs, hostconfig)
		}
	}

	return dbHostConfigs, nil
}

func ConvertDBModelToConfig(dbModels []db.HostConfig) *config.HostConfigYAML {
	yamlConfig := &config.HostConfigYAML{
		Hosts: make(map[string]struct {
			ApiUrl     string `mapstructure:"api-url"`
			TargetNode string `mapstructure:"target-node"`
			LXCs       struct {
				Default  config.MachineConfigYAML   `mapstructure:"default"`
				Machines []config.MachineConfigYAML `mapstructure:"machines"`
			} `mapstructure:"lxc"`
			VMs struct {
				Default  config.MachineConfigYAML   `mapstructure:"default"`
				Machines []config.MachineConfigYAML `mapstructure:"machines"`
			} `mapstructure:"vm"`
		}),
	}

	for _, hostConfig := range dbModels {
		nodeConfig := struct {
			ApiUrl     string `mapstructure:"api-url"`
			TargetNode string `mapstructure:"target-node"`
			LXCs       struct {
				Default  config.MachineConfigYAML   `mapstructure:"default"`
				Machines []config.MachineConfigYAML `mapstructure:"machines"`
			} `mapstructure:"lxc"`
			VMs struct {
				Default  config.MachineConfigYAML   `mapstructure:"default"`
				Machines []config.MachineConfigYAML `mapstructure:"machines"`
			} `mapstructure:"vm"`
		}{
			ApiUrl:     hostConfig.ApiURL,
			TargetNode: hostConfig.TargetNode,
		}

		for _, machine := range hostConfig.Machines {
			machinePubKeys := []config.SSHKey{}
			for _, pubkey := range machine.SSHPublicKeys {
				machinePubKeys = append(machinePubKeys, config.SSHKey{
					Key:  *pubkey.Key,
					Path: *pubkey.Path,
				})
			}

			machineConfig := config.MachineConfigYAML{
				Name:             machine.Name,
				Type:             machine.Type,
				VmId:             machine.VmId,
				OsTemplate:       machine.OsTemplate,
				ISO:              machine.ISO,
				IPAddress:        machine.IPAddress,
				CPUs:             machine.CPUs,
				Memory:           machine.Memory,
				DiskSize:         machine.DiskSize,
				SwapSize:         machine.SwapSize,
				Tags:             machine.Tags,
				StartOnBoot:      machine.StartOnBoot,
				StorageBackend:   machine.StorageBackend,
				TemplateBackend:  machine.TemplateBackend,
				NetworkBridge:    machine.NetworkBridge,
				NetworkInterface: machine.NetworkInterface,
				DefaultGateway:   machine.DefaultGateway,
				SshPublicKeys:    machinePubKeys,
			}

			if machine.Type == "lxc" {
				nodeConfig.LXCs.Machines = append(nodeConfig.LXCs.Machines, machineConfig)
			} else if machine.Type == "vm" {
				nodeConfig.VMs.Machines = append(nodeConfig.VMs.Machines, machineConfig)
			}
		}

		yamlConfig.Hosts[hostConfig.TargetNode] = nodeConfig
	}

	return yamlConfig
}

func PropagateDefaults(cfg *config.HostConfigYAML) *config.HostConfigYAML {
	for _, host := range cfg.Hosts {
		for i, m := range host.LXCs.Machines {
			host.LXCs.Machines[i] = applyDefaults(m, host.LXCs.Default)
		}
		for i, m := range host.VMs.Machines {
			host.VMs.Machines[i] = applyDefaults(m, host.VMs.Default)
		}
	}

	return cfg
}

func applyDefaults(machine config.MachineConfigYAML, defaultConfig config.MachineConfigYAML) config.MachineConfigYAML {
	if machine.NetworkBridge == "" {
		machine.NetworkBridge = defaultConfig.NetworkBridge
	}
	if machine.NetworkInterface == "" {
		machine.NetworkInterface = defaultConfig.NetworkInterface
	}
	if machine.DefaultGateway == "" {
		machine.DefaultGateway = defaultConfig.DefaultGateway
	}
	if machine.StorageBackend == "" {
		machine.StorageBackend = defaultConfig.StorageBackend
	}
	if machine.TemplateBackend == "" {
		machine.TemplateBackend = defaultConfig.TemplateBackend
	}
	if machine.DiskSize == 0 {
		machine.DiskSize = defaultConfig.DiskSize
	}
	if machine.SwapSize == 0 {
		machine.SwapSize = defaultConfig.SwapSize
	}
	if machine.CPUs == 0 {
		machine.CPUs = defaultConfig.CPUs
	}
	if machine.Memory == 0 {
		machine.Memory = defaultConfig.Memory
	}
	if !machine.StartOnBoot {
		machine.StartOnBoot = defaultConfig.StartOnBoot
	}
	if machine.Tags == "" {
		machine.Tags = defaultConfig.Tags
	}
	return machine
}
