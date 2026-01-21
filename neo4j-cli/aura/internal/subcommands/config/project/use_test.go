package project_test

import (
	"testing"

	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
)

func TestUseProject(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.SetConfigValue("projects.projects", []map[string]string{{"name": "test", "organization-id": "testorganizationid", "project-id": "testprojectid"}})

	helper.ExecuteCommand("config project use test")

	helper.AssertConfigValue("projects.default-project", "test")
}

func TestUseProjectIfDoesNotExist(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	helper.SetConfigValue("aura.beta-enabled", true)
	helper.ExecuteCommand("config project use test")

	helper.AssertErr("Error: could not find a project with the name test")
}
