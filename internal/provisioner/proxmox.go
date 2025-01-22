package provisioner

import (
	"fmt"
	"strings"

	"dario.cat/mergo"
	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tscrond/sprinkle/config"
	"github.com/tscrond/sprinkle/internal/auth"
	"github.com/tscrond/sprinkle/internal/configmapper"
	"github.com/tscrond/sprinkle/internal/db"
)

type ProxmoxProvisioner struct {
	AuthService   *auth.AuthService
	ConfigManager *config.ConfigManager
}

func NewProxmoxProvisioner(authService *auth.AuthService, configManager *config.ConfigManager) *ProxmoxProvisioner {
	return &ProxmoxProvisioner{
		AuthService:   authService,
		ConfigManager: configManager,
	}
}

func (p *ProxmoxProvisioner) DestroyResource() error {
	return nil
}

// apply diff flow:
// 0. read YAML config
// 1. read current DB state
// 2. convert DB state object to YAML object (or vice-versa)
// 3. check diff (using some library), save it to variable
// 4. try to apply the diffed config (current desired state)
// 5. if no failure, update database
func (p *ProxmoxProvisioner) ApplyDiff() error {

	newDbState, err := p.ComputeDiff()
	if err != nil {
		return err
	}

	if !p.PromptForApply() {
		fmt.Println("Not applying new state!")
		return nil
	}

	fmt.Println("Applying new state to Proxmox Nodes!")
	// TODO: call proxmox API to provision the resources/apply the changes

	if err := p.AuthService.Db.InsertHostConfigs(newDbState); err != nil {
		fmt.Println("error inserting new state to db ", err)
		return err
	}

	return nil
}

func (p *ProxmoxProvisioner) PromptForApply() bool {
	var apply string
	fmt.Print("Do you want to apply? (yes/no): ")
	fmt.Scanln(&apply)

	doApply := false

	if apply == "yes" {
		doApply = true
	}

	return doApply
}

func (p *ProxmoxProvisioner) ComputeDiff() ([]db.HostConfig, error) {
	stateFromYaml, err := p.ConfigManager.LoadConfigFromYAML()
	if err != nil {
		return nil, err
	}
	stateFromDb, err := p.AuthService.Db.GetAllHostConfigs()
	if err != nil {
		return nil, err
	}
	stateFromDbConf, err := configmapper.ConvertDBModelToConfig(stateFromDb)
	if err != nil {
		return nil, err
	}

	p.DisplayDiff(stateFromYaml, stateFromDbConf)

	if err := mergo.Merge(stateFromDbConf, *stateFromYaml, mergo.WithOverride); err != nil {
		fmt.Println(err)
	}

	mergedNewDbState, err := configmapper.ConvertConfigToDBModel(stateFromDbConf)
	if err != nil {
		return nil, err
	}

	return mergedNewDbState, nil
}

func (p *ProxmoxProvisioner) DisplayDiff(configFromYaml, configFromDB *config.HostConfigYAML) error {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	diff := cmp.Diff(configFromDB, configFromYaml, cmp.AllowUnexported(config.MachineConfigYAML{}), cmpopts.EquateEmpty())

	var result strings.Builder
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") {
			result.WriteString(green(line) + "\n")
		} else if strings.HasPrefix(line, "-") {
			result.WriteString(red(line) + "\n")
		} else {
			result.WriteString(line + "\n")
		}
	}

	fmt.Print(result.String())

	return nil
}
