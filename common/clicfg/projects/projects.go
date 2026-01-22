package projects

import (
	"encoding/json"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/neo4j/cli/common/clierr"
	"github.com/spf13/afero"
	"github.com/tidwall/sjson"
)

type AuraConfigProjects struct {
	fs       afero.Fs
	filePath string
}

type AuraProjectConfig struct {
	Aura     any           `json:"aura"`
	Projects *AuraProjects `json:"aura-projects"`
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

func NewAuraConfigProjects(fs afero.Fs, filePath string) *AuraConfigProjects {
	return &AuraConfigProjects{fs: fs, filePath: filePath}
}

func (p *AuraConfigProjects) Add(name string, organizationId string, projectId string) error {
	data := fileutils.ReadFileSafe(p.fs, p.filePath)

	projects, err := p.projects(data)
	if err != nil {
		return err
	}

	for _, project := range projects.Projects {
		if project.Name == name {
			return clierr.NewUsageError("already have a project with the name %s", name)
		}
	}

	projects.Projects = append(projects.Projects, &AuraProject{Name: name, OrganizationId: organizationId, ProjectId: projectId})
	if len(projects.Projects) == 1 {
		projects.DefaultProject = name
	}

	return p.updateProjects(data, projects)
}

func (p *AuraConfigProjects) Remove(name string) error {
	data := fileutils.ReadFileSafe(p.fs, p.filePath)

	projects, err := p.projects(data)
	if err != nil {
		return err
	}

	indexToRemove := -1
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

	return p.updateProjects(data, projects)
}

func (p *AuraConfigProjects) SetDefault(name string) (*AuraProject, error) {
	data := fileutils.ReadFileSafe(p.fs, p.filePath)

	projects, err := p.projects(data)
	if err != nil {
		return nil, err
	}
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

	err = p.updateProjects(data, projects)
	if err != nil {
		return nil, err
	}
	return defaultProject, nil
}

func (p *AuraConfigProjects) Default() (*AuraProject, error) {
	data := fileutils.ReadFileSafe(p.fs, p.filePath)

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

func (p *AuraConfigProjects) projects(data []byte) (*AuraProjects, error) {
	auraProjectConfig := AuraProjectConfig{}
	if err := json.Unmarshal(data, &auraProjectConfig); err != nil {
		return nil, err
	}

	return auraProjectConfig.Projects, nil
}

func (p *AuraConfigProjects) updateProjects(data []byte, projects *AuraProjects) error {
	updateConfig, err := sjson.Set(string(data), "aura-projects", projects)
	if err != nil {
		return err
	}

	fileutils.WriteFile(p.fs, p.filePath, []byte(updateConfig))
	return nil
}
