package project

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewCmd(cfg *clicfg.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage and view project values",
	}

	cmd.AddCommand(NewAddCmd(cfg))
	cmd.AddCommand(NewUseCmd(cfg))
	cmd.AddCommand(NewListCmd(cfg))
	cmd.AddCommand(NewRemoveCmd(cfg))

	return cmd
}
