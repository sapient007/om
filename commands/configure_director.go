package commands

import (
	"encoding/json"
	"fmt"
	"github.com/pivotal-cf/om/interpolate"
	"os"
	"sort"
	"strings"

	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/om/api"
	"gopkg.in/yaml.v2"
)

type ConfigureDirector struct {
	environFunc func() []string
	service     configureDirectorService
	logger      logger
	Options     struct {
		IgnoreVerifierWarnings bool     `long:"ignore-verifier-warnings" description:"option to ignore verifier warnings. NOT RECOMMENDED UNLESS DISABLED IN OPS MANAGER"`
		ConfigFile             string   `short:"c" long:"config" description:"path to yml file containing all config fields (see docs/configure-director/README.md for format)" required:"true"`
		VarsFile               []string `long:"vars-file" description:"Load variables from a YAML file"`
		VarsEnv                []string `long:"vars-env" description:"Load variables from environment variables (e.g.: 'MY' to load MY_var=value)"`
		Vars                   []string `long:"var" short:"v" description:"Load variable from the command line. Format: VAR=VAL"`
		OpsFile                []string `long:"ops-file" description:"YAML operations file"`
	}
}

type VMTypesConfiguration struct {
	CustomTypesOnly bool               `yaml:"custom_only,omitempty" json:"custom_only,omitempty"`
	VMTypes         []api.CreateVMType `yaml:"vm_types,omitempty" json:"vm_types,omitempty"`
}

type directorConfig struct {
	NetworkAssignment       interface{}            `yaml:"network-assignment"`
	AZConfiguration         interface{}            `yaml:"az-configuration"`
	NetworksConfiguration   interface{}            `yaml:"networks-configuration"`
	PropertiesConfiguration interface{}            `yaml:"properties-configuration"`
	IAASConfigurations      interface{}            `yaml:"iaas-configurations"`
	ResourceConfiguration   map[string]interface{} `yaml:"resource-configuration"`
	VMExtensions            interface{}            `yaml:"vmextensions-configuration"`
	VMTypes                 VMTypesConfiguration   `yaml:"vmtypes-configuration"`
	Field                   map[string]interface{} `yaml:",inline"`
}

//go:generate counterfeiter -o ./fakes/configure_director_service.go --fake-name ConfigureDirectorService . configureDirectorService
type configureDirectorService interface {
	CreateCustomVMTypes(api.CreateVMTypes) error
	CreateStagedVMExtension(api.CreateVMExtension) error
	DeleteCustomVMTypes() error
	DeleteVMExtension(name string) error
	GetStagedProductByName(name string) (api.StagedProductsFindOutput, error)
	GetStagedProductJobResourceConfig(string, string) (api.JobProperties, error)
	GetStagedProductManifest(guid string) (manifest string, err error)
	Info() (api.Info, error)
	ListInstallations() ([]api.InstallationsServiceOutput, error)
	ListStagedProductJobs(string) (map[string]string, error)
	ListStagedVMExtensions() ([]api.VMExtension, error)
	ListVMTypes() ([]api.VMType, error)
	UpdateStagedDirectorIAASConfigurations(api.IAASConfigurationsInput) error
	UpdateStagedDirectorAvailabilityZones(api.AvailabilityZoneInput, bool) error
	UpdateStagedDirectorNetworkAndAZ(api.NetworkAndAZConfiguration) error
	UpdateStagedDirectorNetworks(api.NetworkInput) error
	UpdateStagedDirectorProperties(api.DirectorProperties) error
	UpdateStagedProductJobResourceConfig(string, string, api.JobProperties) error
}

func NewConfigureDirector(environFunc func() []string, service configureDirectorService, logger logger) ConfigureDirector {
	return ConfigureDirector{
		environFunc: environFunc,
		service:     service,
		logger:      logger,
	}
}

func (c ConfigureDirector) Usage() jhanda.Usage {
	return jhanda.Usage{
		Description:      "This authenticated command configures the director.",
		ShortDescription: "configures the director",
		Flags:            c.Options,
	}
}

func (c ConfigureDirector) Execute(args []string) error {
	if _, err := jhanda.Parse(&c.Options, args); err != nil {
		return fmt.Errorf("could not parse configure-director flags: %s", err)
	}

	err := checkRunningInstallation(c.service.ListInstallations)
	if err != nil {
		return err
	}

	config, err := c.interpolateConfig()
	if err != nil {
		return err
	}

	err = c.validateConfig(config)
	if err != nil {
		return err
	}

	err = c.updateIAASConfigurations(config)
	if err != nil {
		return err
	}

	err = c.updateStagedDirectorProperties(config)
	if err != nil {
		return err
	}

	err = c.configureAvailabilityZones(config)
	if err != nil {
		return err
	}

	err = c.configureNetworksConfiguration(config)
	if err != nil {
		return err
	}

	err = c.configureNetworkAssignment(config)
	if err != nil {
		return err
	}

	err = c.configureVMTypes(config)
	if err != nil {
		return err
	}

	err = c.configureResourceConfigurations(config)
	if err != nil {
		return err
	}

	err = c.configureVMExtensions(config)
	if err != nil {
		return err
	}

	return nil
}

func (c ConfigureDirector) interpolateConfig() (*directorConfig, error) {
	varsEnvs := c.Options.VarsEnv
	if value, ok := os.LookupEnv("OM_VARS_ENV"); ok {
		// EXPERIMENTAL: don't put this directly in VarsEnv
		varsEnvs = append(varsEnvs, value)
	}

	configContents, err := interpolate.Execute(interpolate.Options{
		TemplateFile:  c.Options.ConfigFile,
		VarsFiles:     c.Options.VarsFile,
		EnvironFunc:   c.environFunc,
		Vars:          c.Options.Vars,
		VarsEnvs:      varsEnvs,
		OpsFiles:      c.Options.OpsFile,
		ExpectAllKeys: true,
	}, "")
	if err != nil {
		return nil, err
	}

	var config directorConfig
	err = yaml.UnmarshalStrict(configContents, &config)
	if err != nil {
		return nil, fmt.Errorf("could not be parsed as valid configuration: %s: %s", c.Options.ConfigFile, err)
	}
	return &config, nil
}

func (c ConfigureDirector) validateConfig(config *directorConfig) error {
	err := c.checkForDeprecatedKeys(config)
	if err != nil {
		return err
	}

	err = c.checkIAASConfigurationIsOnlySetOnce(config)
	if err != nil {
		return err
	}
	return nil
}

func (c ConfigureDirector) checkForDeprecatedKeys(config *directorConfig) error {
	if len(config.Field) > 0 {
		var unrecognizedKeys []string
		for key := range config.Field {
			unrecognizedKeys = append(unrecognizedKeys, key)
		}
		sort.Strings(unrecognizedKeys)

		deprecatedKeys := []string{"director-configuration", "iaas-configuration", "security-configuration", "syslog-configuration"}
		errorMessage := `The following keys have recently been removed from the top level configuration: director-configuration, iaas-configuration, security-configuration, syslog-configuration
To fix this error, move the above keys under 'properties-configuration' and change their dashes to underscores.

The old configuration file would contain the keys at the top level.

director-configuration: {}
iaas-configuration: {}
network-assignment: {}
networks-configuration: {}
resource-configuration: {}
security-configuration: {}
syslog-configuration: {}
vmextensions-configuration: {}

They'll need to be moved to the new 'properties-configuration', with their dashes turn to underscore.
For example, 'director-configuration' becomes 'director_configuration'.

The new configration file will look like.

az-configuration: {}
network-assignment: {}
networks-configuration: {}
properties-configuration:
  director_configuration: {}
  security_configuration: {}
  syslog_configuration: {}
  iaas_configuration: {}
resource-configuration: {}
vmextensions-configuration: {}
`

		for _, depKey := range deprecatedKeys {
			for _, unrecKey := range unrecognizedKeys {
				if depKey == unrecKey {
					return fmt.Errorf(errorMessage)
				}
			}
		}

		return fmt.Errorf("the config file contains unrecognized keys: \"%s\"", strings.Join(unrecognizedKeys, "\", \""))
	}
	return nil
}

func (c ConfigureDirector) checkIAASConfigurationIsOnlySetOnce(config *directorConfig) error {
	iaasConfigurations := config.IAASConfigurations
	properties, ok := config.PropertiesConfiguration.(map[interface{}]interface{})
	if !ok {
		return nil
	}

	iaasProperties := properties["iaas-configuration"]

	if iaasConfigurations != nil && iaasProperties != nil {
		return fmt.Errorf("iaas-configurations cannot be used with properties-configuration.iaas-configurations\n" +
			"Please only use one implementation.")
	}
	return nil
}

func (c ConfigureDirector) updateIAASConfigurations(config *directorConfig) error {
	if config.IAASConfigurations != nil {
		c.logger.Printf("started setting iaas configurations for bosh tile")

		info, err := c.service.Info()
		if err != nil {
			return fmt.Errorf("could not retrieve info from targetted ops manager: %v", err)
		}
		if ok, _ := info.VersionAtLeast(2, 2); !ok {
			return fmt.Errorf("\"iaas-configurations\" is only available with Ops Manager 2.2 or later: you are running %s", info.Version)
		}

		configurations, err := getJSONProperties(config.IAASConfigurations)
		if err != nil {
			return err
		}

		err = c.service.UpdateStagedDirectorIAASConfigurations(api.IAASConfigurationsInput(configurations))

		if err != nil {
			return fmt.Errorf("iaas configurations could not be completed: %s", err)
		}

		c.logger.Printf("finished setting iaas configurations for bosh tile")
	}
	return nil
}

func (c ConfigureDirector) updateStagedDirectorProperties(config *directorConfig) error {
	if config.PropertiesConfiguration != nil {
		c.logger.Printf("started configuring director options for bosh tile")

		properties, err := getJSONProperties(config.PropertiesConfiguration)
		if err != nil {
			return err
		}

		err = c.service.UpdateStagedDirectorProperties(api.DirectorProperties(properties))

		if err != nil {
			return fmt.Errorf("properties could not be applied: %s", err)
		}

		c.logger.Printf("finished configuring director options for bosh tile")
	}
	return nil
}

func (c ConfigureDirector) configureAvailabilityZones(config *directorConfig) error {
	if config.AZConfiguration != nil {
		c.logger.Printf("started configuring availability zone options for bosh tile")

		azs, err := getJSONProperties(config.AZConfiguration)
		if err != nil {
			return err
		}

		err = c.service.UpdateStagedDirectorAvailabilityZones(api.AvailabilityZoneInput{
			AvailabilityZones: json.RawMessage(azs),
		}, c.Options.IgnoreVerifierWarnings)
		if err != nil {
			return fmt.Errorf("availability zones configuration could not be applied: %s", err)
		}

		c.logger.Printf("finished configuring availability zone options for bosh tile")
	}
	return nil
}

func (c ConfigureDirector) configureNetworksConfiguration(config *directorConfig) error {
	if config.NetworksConfiguration != nil {
		c.logger.Printf("started configuring network options for bosh tile")

		networksConfiguration, err := getJSONProperties(config.NetworksConfiguration)
		if err != nil {
			return err
		}

		err = c.service.UpdateStagedDirectorNetworks(api.NetworkInput{
			Networks: json.RawMessage(networksConfiguration),
		})
		if err != nil {
			return fmt.Errorf("networks configuration could not be applied: %s", err)
		}

		c.logger.Printf("finished configuring network options for bosh tile")
	}
	return nil
}

func (c ConfigureDirector) configureNetworkAssignment(config *directorConfig) error {
	if config.NetworkAssignment != nil {
		c.logger.Printf("started configuring network assignment options for bosh tile")

		networkAssignment, err := getJSONProperties(config.NetworkAssignment)
		if err != nil {
			return err
		}
		err = c.service.UpdateStagedDirectorNetworkAndAZ(api.NetworkAndAZConfiguration{
			NetworkAZ: json.RawMessage(networkAssignment),
		})
		if err != nil {
			return fmt.Errorf("network and AZs could not be applied: %s", err)
		}

		c.logger.Printf("finished configuring network assignment options for bosh tile")
	}

	return nil
}

func (c ConfigureDirector) configureResourceConfigurations(config *directorConfig) error {
	if config.ResourceConfiguration != nil {
		c.logger.Printf("started configuring resource options for bosh tile")

		productGUID, err := c.getProductGUID()
		if err != nil {
			return err
		}

		jobNames := c.getUserProvidedJobNames(config)

		jobs, err := c.service.ListStagedProductJobs(productGUID)
		if err != nil {
			return fmt.Errorf("failed to fetch jobs: %s", err)
		}

		c.logger.Printf("applying resource configuration for the following jobs:")
		for _, jobName := range jobNames {
			err := c.configureResourceConfiguration(config, jobs, jobName, productGUID)
			if err != nil {
				return err
			}
		}

		c.logger.Printf("finished configuring resource options for bosh tile")
	}
	return nil
}

func (c ConfigureDirector) addNewExtensions(extensionsToDelete map[string]api.VMExtension, newExtensions interface{}) (map[string]api.VMExtension, error) {
	var newVMExtensions []api.VMExtension

	newExtensionBytes, err := getJSONProperties(newExtensions)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(newExtensionBytes), &newVMExtensions)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall vmextensions-configuration json: %s. Full Error: %s", newExtensions, err)
	}

	c.logger.Printf("applying vmextensions configuration for the following:")
	for _, newExtension := range newVMExtensions {
		c.logger.Printf("\t%s", newExtension.Name)

		cloudProperties, err := json.Marshal(newExtension.CloudProperties)
		if err != nil {
			return nil, err
		}

		err = c.service.CreateStagedVMExtension(api.CreateVMExtension{
			Name:            newExtension.Name,
			CloudProperties: cloudProperties,
		})
		if err != nil {
			return nil, err
		}

		for name := range extensionsToDelete {
			if name == newExtension.Name {
				delete(extensionsToDelete, name)
			}
		}
	}
	return extensionsToDelete, nil
}

func (c ConfigureDirector) getExistingExtensions() (map[string]api.VMExtension, error) {
	existingVMExtensions, err := c.service.ListStagedVMExtensions()
	if err != nil {
		return nil, err
	}

	extensionsToDelete := make(map[string]api.VMExtension)
	for _, vmExtension := range existingVMExtensions {
		extensionsToDelete[vmExtension.Name] = vmExtension
	}

	return extensionsToDelete, nil
}

func (c ConfigureDirector) deleteExtensions(extensionsToDelete map[string]api.VMExtension) error {
	for _, extensionToDelete := range extensionsToDelete {
		c.logger.Printf("deleting vm extension %s", extensionToDelete.Name)
		err := c.service.DeleteVMExtension(extensionToDelete.Name)
		if err != nil {
			return err
		}
		c.logger.Printf("done deleting vm extension %s", extensionToDelete.Name)
	}
	return nil
}

func (c ConfigureDirector) configureVMExtensions(config *directorConfig) error {
	if config.VMExtensions != nil {
		c.logger.Printf("started configuring vm extensions")

		currentExtensions, err := c.getExistingExtensions()
		if err != nil {
			return err
		}

		extensionsToDelete, err := c.addNewExtensions(currentExtensions, config.VMExtensions)
		if err != nil {
			return err
		}

		err = c.deleteExtensions(extensionsToDelete)
		if err != nil {
			return err
		}

		c.logger.Printf("finished configuring vm extensions")
	}
	return nil
}

func (c ConfigureDirector) configureVMTypes(config *directorConfig) error {
	if len(config.VMTypes.VMTypes) == 0 {
		if config.VMTypes.CustomTypesOnly {
			return fmt.Errorf("if custom_types = true, vm_types must not be empty")
		}

		return nil
	}

	c.logger.Printf("creating custom vm types")

	vmTypesToCreate := make([]api.CreateVMType, 0)
	existingVMTypes := make([]api.VMType, 0)

	var err error
	if !config.VMTypes.CustomTypesOnly {
		// delete all custom VM types
		if err = c.service.DeleteCustomVMTypes(); err != nil {
			return err
		}

		existingVMTypes, err = c.service.ListVMTypes()
		if err != nil {
			return err
		}
	}

	for i := range existingVMTypes {
		vmTypesToCreate = append(vmTypesToCreate, existingVMTypes[i].CreateVMType)
	}

	for i := range config.VMTypes.VMTypes {
		overwriting := false
		for j := range vmTypesToCreate {
			if config.VMTypes.VMTypes[i].Name == vmTypesToCreate[j].Name {
				vmTypesToCreate[j] = config.VMTypes.VMTypes[i]
				overwriting = true
				break
			}
		}

		if !overwriting {
			vmTypesToCreate = append(vmTypesToCreate, config.VMTypes.VMTypes[i])
		}
	}

	return c.service.CreateCustomVMTypes(api.CreateVMTypes{
		VMTypes: vmTypesToCreate,
	})
}

func (c ConfigureDirector) getProductGUID() (string, error) {
	findOutput, err := c.service.GetStagedProductByName("p-bosh")
	if err != nil {
		return "", fmt.Errorf("could not find staged product with name 'p-bosh': %s", err)
	}
	return findOutput.Product.GUID, nil
}

func (c ConfigureDirector) getUserProvidedJobNames(config *directorConfig) []string {
	var names []string

	for name := range config.ResourceConfiguration {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func (c ConfigureDirector) configureResourceConfiguration(config *directorConfig, jobs map[string]string, jobName string, productGUID string) error {
	c.logger.Printf("\t%s", jobName)
	jobGUID, ok := jobs[jobName]
	if !ok {
		return fmt.Errorf("product 'p-bosh' does not contain a job named '%s'", jobName)
	}

	jobProperties, err := c.service.GetStagedProductJobResourceConfig(productGUID, jobGUID)
	if err != nil {
		return fmt.Errorf("could not fetch existing job configuration for '%s': %s", jobName, err)
	}

	prop, err := getJSONProperties(config.ResourceConfiguration[jobName])
	if err != nil {
		return fmt.Errorf("could not unmarshall resource configuration: %v", err)
	}

	err = json.Unmarshal([]byte(prop), &jobProperties)
	if err != nil {
		return fmt.Errorf("could not decode resource-configuration json for job '%s': %s", jobName, err)
	}

	err = c.service.UpdateStagedProductJobResourceConfig(productGUID, jobGUID, jobProperties)
	if err != nil {
		return fmt.Errorf("failed to configure resources for '%s': %s", jobName, err)
	}

	return nil
}

func checkRunningInstallation(listInstallations func() ([]api.InstallationsServiceOutput, error)) error {
	installations, err := listInstallations()
	if err != nil {
		return fmt.Errorf("could not list installations: %s", err)
	}
	if len(installations) > 0 {
		if installations[0].Status == "running" {
			return fmt.Errorf("OpsManager does not allow configuration or staging changes while apply changes are running to prevent data loss for configuration and/or staging changes")
		}
	}
	return nil
}
