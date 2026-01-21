package project_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestAddFirstProject(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("projects.projects", []map[string]string{})

	helper.ExecuteCommand("config project add --name test --organization-id testorganizationid --project-id testprojectid")

	helper.AssertConfigValue("projects.projects", `[{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}]`)
	helper.AssertConfigValue("projects.default-project", "test")
}

func TestAddProjectIfAlreadyExists(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("projects.projects", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})

	helper.ExecuteCommand("config project add --name test --organization-id testorganizationid --project-id testprojectid")

	helper.AssertErr("Error: already have a project with the name test")
}
func TestAddAditionalProjects(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("projects.projects", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})
	helper.SetConfigValue("projects.default-project", "test")

	helper.ExecuteCommand("config project add --name test-new --organization-id newtestorganizationid --project-id newtestprojectid")

	helper.AssertConfigValue("projects.projects", `[{"name":"test","organization-id":"testorganizationid","project-id":"testprojectid"}, {"name":"test-new","organization-id":"newtestorganizationid","project-id":"newtestprojectid"}]`)
	helper.AssertConfigValue("projects.default-project", "test")
}
