package state

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
	"github.com/tscrond/sprinkle/internal/provisioner"
)

type StateEngine struct {
	AuthService   *auth.AuthService
	ConfigManager *config.ConfigManager
	Provisioner   *provisioner.ProxmoxProvisioner
}

func NewStateEngine(authService *auth.AuthService, configManager *config.ConfigManager, provisioner *provisioner.ProxmoxProvisioner) *StateEngine {
	return &StateEngine{
		AuthService:   authService,
		ConfigManager: configManager,
		Provisioner:   provisioner,
	}
}

// apply diff flow:
// 0. read YAML config
// 1. read current DB state
// 2. convert DB state object to YAML object (or vice-versa)
// 3. check diff (using some library), save it to variable
// 4. try to apply the diffed config (current desired state)
// 5. if no failure, update database
func (p *StateEngine) ApplyDiff(targetNode string) error {

	newDbState, err := p.ComputeDiff()
	if err != nil {
		return err
	}

	if !p.PromptForApply(targetNode) {
		fmt.Println("Not applying new state!")
		return nil
	}

	fmt.Println("Applying new state to Proxmox Node!")
	// TODO: call proxmox API to provision the resources/apply the changes

	if err := p.Provisioner.ApplyNewState(targetNode, newDbState); err != nil {
		fmt.Println("error applying state ", err)
	}

	if err := p.AuthService.Db.InsertHostConfigs(newDbState); err != nil {
		fmt.Println("error inserting new state to db ", err)
		return err
	}

	return nil
}

func (p *StateEngine) PromptForApply(targetNode string) bool {
	var apply string
	fmt.Println("WARNING, altering node: ", targetNode)
	fmt.Print("Do you want to apply? (yes/no): ")
	fmt.Scanln(&apply)

	doApply := false

	if apply == "yes" {
		doApply = true
	}

	return doApply
}

func (p *StateEngine) ComputeDiff() ([]db.HostConfig, error) {
	stateFromYaml, err := p.ConfigManager.LoadConfigFromYAML()
	if err != nil {
		return nil, err
	}
	stateFromDb, err := p.AuthService.Db.GetAllHostConfigs()
	if err != nil {
		return nil, err
	}
	stateFromDbConf := configmapper.ConvertDBModelToConfig(stateFromDb)

	stateFromYaml = configmapper.PropagateDefaults(stateFromYaml)

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

func (p *StateEngine) DisplayDiff(configFromYaml, configFromDB *config.HostConfigYAML) error {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	diff := cmp.Diff(configFromDB, configFromYaml, cmpopts.EquateEmpty())

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
