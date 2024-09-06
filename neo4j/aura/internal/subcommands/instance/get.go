package instance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Returns instance details",
		Long:  "This endpoint returns details about a specific Aura Instance.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			aura, err := api.GetApiFromConfig(cmd)
			if err != nil {
				return fmt.Errorf("error in command %s", os.Args[1:])
			}

			instanceId := args[0]
			instances, statusCode, err := aura.Instances.Get(instanceId)
			if err != nil {
				return err
			}

			if statusCode == http.StatusOK {
				jsonResponse, err := json.Marshal(instances)
				if err != nil {
					return fmt.Errorf("error in command %s: %v", os.Args[1:], err)
				}

				err = output.PrintBody(cmd, jsonResponse, []string{"id", "name", "tenant_id", "status", "connection_url", "cloud_provider", "region", "type", "memory", "storage", "customer_managed_key_id"})
				if err != nil {
					return fmt.Errorf("error in command %s: %v", os.Args[1:], err)
				}
			}

			return nil

			// path := fmt.Sprintf("/instances/%s", args[0])

			// resBody, statusCode, err := api.MakeRequest(cmd, http.MethodGet, path, nil)
			// if err != nil {
			// 	return err
			// }

			// if statusCode == http.StatusOK {
			// 	err = output.PrintBody(cmd, resBody)
			// 	if err != nil {
			// 		return err
			// 	}

			// }

			// return nil
		},
	}
}
