package utils

import (
	"log"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

// This function is meant to run in the PreRun of V2 commands to ensure that the flags are marked as required if no values have been set
// through the `config project add/use` commands.
func SetOragnizationAndProjectIdFlagsAsRequired(cfg *clicfg.Config, cmd *cobra.Command, organizationIdFlag string, projectIdFlag string) {
	defaultProject, err := cfg.Aura.Projects.Default()
	if err != nil {
		log.Fatal(err)
	}
	if defaultProject.OrganizationId == "" {
		err := cmd.MarkFlagRequired(organizationIdFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	if defaultProject.ProjectId == "" {
		err := cmd.MarkFlagRequired(projectIdFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// This function is meant to run in the RunE of V2 commands to ensure that the values are set as the given default values if no values are
// given via flags when running the command.
func SetMissingOragnizationAndProjectIdValuesFromDefaults(cfg *clicfg.Config, organizationId *string, projectId *string) {
	defaultProject, err := cfg.Aura.Projects.Default()
	if err != nil {
		log.Fatal(err)
	}
	if projectId != nil && *organizationId == "" {
		*organizationId = defaultProject.OrganizationId
	}
	if projectId != nil && *projectId == "" {
		*projectId = defaultProject.ProjectId
	}
}
