package project_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestListProjects(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("output", "json")
	helper.SetConfigValue("projects.projects", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})
	helper.SetConfigValue("projects.default-project", "test")

	helper.ExecuteCommand("config project list")

	helper.AssertOutJson(`
		{
			"default-project": "test",
			"projects": [
				{
					"name": "test",
					"organization-id": "testorganizationid",
					"project-id": "testprojectid"
				}
			]
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
		"projects": []
	}`)
}
