package subcommands

import (
	"log"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func SetFlagsAsRequired(cfg *clicfg.Config, cmd *cobra.Command, organizationIdFlag string, projectIdFlag string) error {
	defaultProject, err := cfg.Aura.GetDefaultProject()
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

	return nil
}

func SetMissingValuesFromDefaults(cfg *clicfg.Config, organizationId *string, projectId *string) {
	defaultProject, err := cfg.Aura.GetDefaultProject()
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
