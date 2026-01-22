package project_test

import (
	"testing"

	"github.com/neo4j/cli/common/clicfg/projects"
	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestListProjects(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("output", "json")
	helper.SetConfigValue("aura-projects.projects", map[string]*projects.AuraProject{"test": {OrganizationId: "testorganizationid", ProjectId: "testprojectid"}})
	helper.SetConfigValue("aura-projects.default-project", "test")

	helper.ExecuteCommand("config project list")

	helper.AssertOutJson(`
		{
			"default-project": "test",
			"projects": {
				"test": {
					"organization-id": "testorganizationid",
					"project-id": "testprojectid"
				}
			}
		}
	`)
}

func TestListProjectWithNoData(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("output", "json")
	helper.ExecuteCommand("config project list")

	helper.AssertOutJson(`{
		"default-project": "",
		"projects": {}
	}`)
}
