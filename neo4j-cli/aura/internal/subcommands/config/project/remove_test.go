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
