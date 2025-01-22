package configmapper

import (
	"github.com/jinzhu/copier"

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
	var dbHostConfigs []db.HostConfig

	for _, hostConfigYaml := range cfg.Hosts {
		var dbHostConf db.HostConfig
		if err := copier.Copy(&dbHostConf, &hostConfigYaml); err != nil {
			return nil, err
		}
		// Map LXCs
		var lxcConfigs []db.LXCConfig
		for _, machine := range hostConfigYaml.LXCs.Machines {
			var lxc db.LXCConfig
			var dbMachines []db.MachineConfig

			if err := copier.Copy(&dbMachines, &machine); err != nil {
				return nil, err
			}
			lxc.Machines = dbMachines
			lxcConfigs = append(lxcConfigs, lxc)
		}
		dbHostConf.LXCs = lxcConfigs

		// Map VMs
		var vmConfigs []db.VMConfig
		for _, machine := range hostConfigYaml.VMs.Machines {
			var vm db.VMConfig
			var dbMachines []db.MachineConfig

			if err := copier.Copy(&dbMachines, &machine); err != nil {
				return nil, err
			}
			vm.Machines = dbMachines
			vmConfigs = append(vmConfigs, vm)
		}
		dbHostConf.VMs = vmConfigs

		dbHostConfigs = append(dbHostConfigs, dbHostConf)
	}
	return dbHostConfigs, nil
}

func ConvertDBModelToConfig(dbHostConfigs []db.HostConfig) (*config.HostConfigYAML, error) {
	var cfg config.HostConfigYAML
	cfg.Hosts = make(map[string]struct {
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
	})

	// Iterate over the DB host configurations
	for _, dbHostConfig := range dbHostConfigs {
		hostConfigYaml := struct {
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
		}{}

		// Map HostConfig fields (ApiUrl, TargetNode)
		hostConfigYaml.ApiUrl = dbHostConfig.ApiURL
		hostConfigYaml.TargetNode = dbHostConfig.TargetNode

		// Map LXCs
		for _, lxcConfig := range dbHostConfig.LXCs {
			// Map LXC Machines
			for _, machineConfig := range lxcConfig.Machines {
				machineYaml := config.MachineConfigYAML{
					Name:             machineConfig.Name,
					VmId:             machineConfig.VmId,
					OsTemplate:       machineConfig.OsTemplate,
					ISO:              machineConfig.ISO,
					IPAddress:        machineConfig.IPAddress,
					CPUs:             machineConfig.CPUs,
					Memory:           machineConfig.Memory,
					DiskSize:         machineConfig.DiskSize,
					SwapSize:         machineConfig.SwapSize,
					Tags:             machineConfig.Tags,
					StartOnBoot:      machineConfig.StartOnBoot,
					StorageBackend:   machineConfig.StorageBackend,
					TemplateBackend:  machineConfig.TemplateBackend,
					NetworkBridge:    machineConfig.NetworkBridge,
					NetworkInterface: machineConfig.NetworkInterface,
					DefaultGateway:   machineConfig.DefaultGateway,
				}

				// Map SSH keys
				var sshKeys []config.SSHKey
				for _, dbSSHKey := range machineConfig.SSHPublicKeys {
					sshKey := config.SSHKey{
						Path: dbSSHKey.Path,
						Key:  dbSSHKey.Key,
					}
					sshKeys = append(sshKeys, sshKey)
				}
				machineYaml.SshPublicKeys = sshKeys

				hostConfigYaml.LXCs.Machines = append(hostConfigYaml.LXCs.Machines, machineYaml)
			}
		}

		// Map VMs
		for _, vmConfig := range dbHostConfig.VMs {
			// Map VM Machines
			for _, machineConfig := range vmConfig.Machines {
				machineYaml := config.MachineConfigYAML{
					Name:             machineConfig.Name,
					VmId:             machineConfig.VmId,
					OsTemplate:       machineConfig.OsTemplate,
					ISO:              machineConfig.ISO,
					IPAddress:        machineConfig.IPAddress,
					CPUs:             machineConfig.CPUs,
					Memory:           machineConfig.Memory,
					DiskSize:         machineConfig.DiskSize,
					SwapSize:         machineConfig.SwapSize,
					Tags:             machineConfig.Tags,
					StartOnBoot:      machineConfig.StartOnBoot,
					StorageBackend:   machineConfig.StorageBackend,
					TemplateBackend:  machineConfig.TemplateBackend,
					NetworkBridge:    machineConfig.NetworkBridge,
					NetworkInterface: machineConfig.NetworkInterface,
					DefaultGateway:   machineConfig.DefaultGateway,
				}

				// Map SSH keys
				var sshKeys []config.SSHKey
				for _, dbSSHKey := range machineConfig.SSHPublicKeys {
					sshKey := config.SSHKey{
						Path: dbSSHKey.Path,
						Key:  dbSSHKey.Key,
					}
					sshKeys = append(sshKeys, sshKey)
				}
				machineYaml.SshPublicKeys = sshKeys

				hostConfigYaml.VMs.Machines = append(hostConfigYaml.VMs.Machines, machineYaml)
			}
		}

		// Add the host config to the final map
		cfg.Hosts[dbHostConfig.TargetNode] = hostConfigYaml
	}

	return &cfg, nil
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
