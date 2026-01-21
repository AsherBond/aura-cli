package project

import (
	"github.com/neo4j/cli/common/clicfg"
	"github.com/spf13/cobra"
)

func NewListCmd(cfg *clicfg.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.Aura.ListProjects(cmd)
		},
	}
}
