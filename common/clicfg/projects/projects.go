package projects

import (
	"encoding/json"
	"fmt"

	"github.com/neo4j/cli/common/clicfg/fileutils"
	"github.com/neo4j/cli/common/clierr"
	"github.com/spf13/afero"
	"github.com/tidwall/sjson"
)

type AuraConfigProjects struct {
	fs       afero.Fs
	filePath string
}

type ConfigAuraProjects struct {
	Projects *AuraProjects `json:"aura-projects"`
}

type AuraProjects struct {
	DefaultProject string                  `json:"default-project"`
	Projects       map[string]*AuraProject `json:"projects"`
}

type AuraProject struct {
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

	if projects == nil {
		projects = &AuraProjects{
			DefaultProject: "",
			Projects:       map[string]*AuraProject{},
		}
	}

	if _, ok := projects.Projects[name]; ok {
		return clierr.NewUsageError("already have a project with the name %s", name)
	}

	projects.Projects[name] = &AuraProject{OrganizationId: organizationId, ProjectId: projectId}

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

	if _, ok := projects.Projects[name]; ok {
		delete(projects.Projects, name)
		if len(projects.Projects) == 0 {
			projects.DefaultProject = ""
		} else {
			if _, ok := projects.Projects[projects.DefaultProject]; !ok {
				for key := range projects.Projects {
					fmt.Printf("Removed the current default project %s, setting %s as the new default project", name, key)
					projects.DefaultProject = key
					break
				}
			}
		}
		return p.updateProjects(data, projects)
	}

	return clierr.NewUsageError("could not find a project with the name %s to remove", name)
}

func (p *AuraConfigProjects) SetDefault(name string) (*AuraProject, error) {
	data := fileutils.ReadFileSafe(p.fs, p.filePath)

	projects, err := p.projects(data)
	if err != nil {
		return nil, err
	}

	if project, ok := projects.Projects[name]; ok {
		projects.DefaultProject = name

		err = p.updateProjects(data, projects)
		if err != nil {
			return nil, err
		}
		return project, nil
	}

	return nil, clierr.NewUsageError("could not find a project with the name %s", name)
}

func (p *AuraConfigProjects) Default() (*AuraProject, error) {
	data := fileutils.ReadFileSafe(p.fs, p.filePath)

	projects, err := p.projects(data)
	if err != nil {
		return nil, err
	}

	if project, ok := projects.Projects[projects.DefaultProject]; ok {
		return project, nil
	}

	return &AuraProject{}, nil
}

func (p *AuraConfigProjects) projects(data []byte) (*AuraProjects, error) {
	auraProjectConfig := ConfigAuraProjects{}
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
