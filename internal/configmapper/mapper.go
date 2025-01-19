package configmapper

import (
	"github.com/jinzhu/copier"
	"github.com/tscrond/sprinkle/config"
	"github.com/tscrond/sprinkle/internal/db"
)

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
