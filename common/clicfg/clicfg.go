package clicfg

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"slices"

	"github.com/neo4j/cli/common/clicfg/credentials"
	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/neo4j/cli/common/clierr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

var ConfigPrefix string

const (
	DefaultAuraBaseUrl     = "https://api.neo4j.io"
	DefaultAuraAuthUrl     = "https://api.neo4j.io/oauth/token"
	DefaultAuraBetaEnabled = false
)

var ValidOutputValues = [3]string{"default", "json", "table"}

type Config struct {
	Version     string
	Aura        *AuraConfig
	Credentials *credentials.Credentials
}

func NewConfig(fs afero.Fs, version string) *Config {
	configPath := filepath.Join(ConfigPrefix, "neo4j", "cli")
	fullConfigPath := filepath.Join(configPath, "config.json")

	Viper := viper.New()

	Viper.SetFs(fs)
	Viper.SetConfigName("config")
	Viper.SetConfigType("json")
	Viper.AddConfigPath(configPath)
	Viper.SetConfigPermissions(0600)

	bindEnvironmentVariables(Viper)
	setDefaultValues(Viper)

	if !fileutils.FileExists(fs, fullConfigPath) {
		if err := fs.MkdirAll(configPath, 0755); err != nil {
			panic(err)
		}
		if err := Viper.SafeWriteConfig(); err != nil {
			panic(err)
		}
	}

	if err := Viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot read config file.")
		panic(err)
	}

	credentials := credentials.NewCredentials(fs, ConfigPrefix)

	return &Config{
		Version: version,
		Aura: &AuraConfig{
			fs:    fs,
			viper: Viper, pollingOverride: PollingConfig{
				MaxRetries: 60,
				Interval:   20,
			},
			ValidConfigKeys: []string{"auth-url", "base-url", "default-tenant", "output", "beta-enabled"},
		},
		Credentials: credentials,
	}
}

func bindEnvironmentVariables(Viper *viper.Viper) {
	Viper.BindEnv("aura.base-url", "AURA_BASE_URL")
	Viper.BindEnv("aura.auth-url", "AURA_AUTH_URL")
}

func setDefaultValues(Viper *viper.Viper) {
	Viper.SetDefault("aura.base-url", DefaultAuraBaseUrl)
	Viper.SetDefault("aura.auth-url", DefaultAuraAuthUrl)
	Viper.SetDefault("aura.output", "default")
	Viper.SetDefault("aura.beta-enabled", DefaultAuraBetaEnabled)
	Viper.SetDefault("projects", AuraProjects{DefaultProject: "", Projects: []*AuraProject{}})
}

type AuraConfig struct {
	viper           *viper.Viper
	fs              afero.Fs
	pollingOverride PollingConfig
	ValidConfigKeys []string
}

type PollingConfig struct {
	Interval   int
	MaxRetries int
}

func (config *AuraConfig) IsValidConfigKey(key string) bool {
	return slices.Contains(config.ValidConfigKeys, key)
}

func (config *AuraConfig) Get(key string) interface{} {
	return config.viper.Get(fmt.Sprintf("aura.%s", key))
}

func (config *AuraConfig) Set(key string, value string) {
	filename := config.viper.ConfigFileUsed()
	data := fileutils.ReadFileSafe(config.fs, filename)

	updateConfig, err := sjson.Set(string(data), fmt.Sprintf("aura.%s", key), value)
	if err != nil {
		panic(err)
	}

	if key == "base-url" {
		updatedAuraBaseUrl := config.auraBaseUrlOnConfigChange(value)
		intermediateUpdateConfig, err := sjson.Set(string(updateConfig), "aura.base-url", updatedAuraBaseUrl)
		if err != nil {
			panic(err)
		}
		updateConfig = intermediateUpdateConfig
	}

	fileutils.WriteFile(config.fs, filename, []byte(updateConfig))
}

func (config *AuraConfig) Print(cmd *cobra.Command) {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config.viper.Get("aura")); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) BaseUrl() string {
	originalUrl := config.viper.GetString("aura.base-url")
	//Existing users have base url configs with trailing path /v1.
	//To make it backward compatible, we allow old config and clear up by removing trailing path /v1 in the url
	return removePathParametersFromUrl(originalUrl)
}

func removePathParametersFromUrl(originalUrl string) string {
	parsedUrl, err := url.Parse(originalUrl)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
}

func (config *AuraConfig) BetaPathV1() string {
	return "v1beta5"
}

func (config *AuraConfig) BetaPathV2() string {
	return "v2beta1"
}

func (config *AuraConfig) BindBaseUrl(flag *pflag.Flag) {
	if err := config.viper.BindPFlag("aura.base-url", flag); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) AuthUrl() string {
	return config.viper.GetString("aura.auth-url")
}

func (config *AuraConfig) BindAuthUrl(flag *pflag.Flag) {
	if err := config.viper.BindPFlag("aura.auth-url", flag); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) Output() string {
	return config.viper.GetString("aura.output")
}

func (config *AuraConfig) BindOutput(flag *pflag.Flag) {
	if err := config.viper.BindPFlag("aura.output", flag); err != nil {
		panic(err)
	}
}

func (config *AuraConfig) AuraBetaEnabled() bool {
	return config.viper.GetBool("aura.beta-enabled")
}

func (config *AuraConfig) DefaultTenant() string {
	return config.viper.GetString("aura.default-tenant")
}

func (config *AuraConfig) Fs() afero.Fs {
	return config.fs
}

func (config *AuraConfig) PollingConfig() PollingConfig {
	return config.pollingOverride
}

func (config *AuraConfig) SetPollingConfig(maxRetries int, interval int) {
	config.pollingOverride = PollingConfig{
		MaxRetries: maxRetries,
		Interval:   interval,
	}
}

func (config *AuraConfig) auraBaseUrlOnConfigChange(url string) string {
	if url == "" {
		return DefaultAuraBaseUrl
	}
	return removePathParametersFromUrl(url)
}

// Types and functions for the "projects" part of the Aura Config
type AuraProjectConfig struct {
	Aura     any           `json:"aura"`
	Projects *AuraProjects `json:"projects"`
}

type AuraProjects struct {
	DefaultProject string         `json:"default-project"`
	Projects       []*AuraProject `json:"projects"`
}

type AuraProject struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organization-id"`
	ProjectId      string `json:"project-id"`
}

func (config *AuraConfig) AddProject(name string, organizationId string, projectId string) error {
	filename := config.viper.ConfigFileUsed()
	data := fileutils.ReadFileSafe(config.fs, filename)

	auraProjectConfig := AuraProjectConfig{Projects: &AuraProjects{DefaultProject: "", Projects: []*AuraProject{}}}
	if err := json.Unmarshal(data, &auraProjectConfig); err != nil {
		return err
	}

	projects := auraProjectConfig.Projects
	for _, project := range projects.Projects {
		if project.Name == name {
			return clierr.NewUsageError("already have a project with the name %s", name)
		}
	}

	projects.Projects = append(projects.Projects, &AuraProject{Name: name, OrganizationId: organizationId, ProjectId: projectId})
	if len(projects.Projects) == 1 {
		projects.DefaultProject = name
	}

	updateConfig, err := sjson.Set(string(data), "projects", projects)
	if err != nil {
		return err
	}

	fileutils.WriteFile(config.fs, filename, []byte(updateConfig))
	return nil
}

func (config *AuraConfig) RemoveProject(name string) error {
	filename := config.viper.ConfigFileUsed()
	data := fileutils.ReadFileSafe(config.fs, filename)

	auraProjectConfig := AuraProjectConfig{}
	if err := json.Unmarshal(data, &auraProjectConfig); err != nil {
		return err
	}
	indexToRemove := -1
	projects := auraProjectConfig.Projects
	for i, project := range projects.Projects {
		if project.Name == name {
			indexToRemove = i
		}
	}

	if indexToRemove == -1 {
		return clierr.NewUsageError("could not find a project with the name %s to remove", name)
	}

	projects.Projects = append(projects.Projects[:indexToRemove], projects.Projects[indexToRemove+1:]...)
	if len(projects.Projects) == 0 {
		projects.DefaultProject = ""
	}

	updateConfig, err := sjson.Set(string(data), "projects", projects)
	if err != nil {
		return err
	}

	fileutils.WriteFile(config.fs, filename, []byte(updateConfig))
	return nil
}

func (config *AuraConfig) ListProjects(cmd *cobra.Command) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config.viper.Get("projects")); err != nil {
		return err
	}
	return nil
}

func (config *AuraConfig) SetDefaultProject(name string) (*AuraProject, error) {
	filename := config.viper.ConfigFileUsed()
	data := fileutils.ReadFileSafe(config.fs, filename)

	auraProjectConfig := AuraProjectConfig{}
	if err := json.Unmarshal(data, &auraProjectConfig); err != nil {
		return nil, err
	}

	projects := auraProjectConfig.Projects
	var defaultProject *AuraProject
	for _, project := range projects.Projects {
		if project.Name == name {
			defaultProject = project
		}
	}
	if defaultProject == nil {
		return nil, clierr.NewUsageError("could not find a project with the name %s", name)
	}
	projects.DefaultProject = name

	updateConfig, err := sjson.Set(string(data), "projects", projects)
	if err != nil {
		return nil, err
	}

	fileutils.WriteFile(config.fs, filename, []byte(updateConfig))
	return defaultProject, nil
}

func (config *AuraConfig) GetDefaultProject() (*AuraProject, error) {
	filename := config.viper.ConfigFileUsed()
	data := fileutils.ReadFileSafe(config.fs, filename)

	auraProjectConfig := AuraProjectConfig{}
	if err := json.Unmarshal(data, &auraProjectConfig); err != nil {
		return nil, err
	}

	projects := auraProjectConfig.Projects
	for _, project := range projects.Projects {
		if project.Name == projects.DefaultProject {
			return project, nil
		}
	}

	return &AuraProject{}, nil
}
