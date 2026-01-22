package project_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestRemoveProject(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura-projects.projects", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})

	helper.ExecuteCommand("config project remove test")

	helper.AssertConfigValue("aura-projects.projects", "[]")
	helper.AssertConfigValue("aura-projects.default-project", "")
}

func TestRemoveProjectWhenProjectDoesNotExist(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura-projects.projects", []map[string]string{})

	helper.ExecuteCommand("config project remove test")

	helper.AssertErr("Error: could not find a project with the name test to remove")
}

func TestRemoveProjectWhenMultipleProjectsExist(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura-projects.projects", []map[string]string{
		{"name": "first-project", "organization-id": "testorganizationid", "project-id": "testprojectid"},
		{"name": "second-project", "organization-id": "testorganizationid", "project-id": "testprojectid"},
	})
	helper.SetConfigValue("aura-projects.default-project", "first-project")

	helper.ExecuteCommand("config project remove first-project")

	helper.AssertConfigValue("aura-projects.projects", `[{"name": "second-project", "organization-id": "testorganizationid", "project-id": "testprojectid"}]`)
	helper.AssertConfigValue("aura-projects.default-project", "second-project")
}

func TestRemoveProjectWhenProjectDoesNotExistWithMultipleProjects(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("aura-projects.projects", []map[string]string{
		{"name": "first-project", "organization-id": "testorganizationid", "project-id": "testprojectid"},
		{"name": "second-project", "organization-id": "testorganizationid", "project-id": "testprojectid"},
	})
	helper.SetConfigValue("aura-projects.default-project", "first-project")

	helper.ExecuteCommand("config project remove non-existing")

	helper.AssertErr("Error: could not find a project with the name non-existing to remove")
	helper.AssertConfigValue("aura-projects.projects", `[{"name": "first-project", "organization-id": "testorganizationid", "project-id": "testprojectid"},{"name": "second-project", "organization-id": "testorganizationid", "project-id": "testprojectid"}]`)
	helper.AssertConfigValue("aura-projects.default-project", "first-project")
}
