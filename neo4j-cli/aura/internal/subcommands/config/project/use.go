package project

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewUseCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "use <name>",
		Short: "Sets the default project to be used",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			defaultProject, err := cfg.Aura.SetDefaultProject(args[0])
			if err != nil {
				return err
			}
			cmd.Printf("Set %s as default project with organization ID %s and project ID %s",
				defaultProject.Name,
				defaultProject.OrganizationId,
				defaultProject.ProjectId,
			)
			return nil
		},
	}
}
