package instance

import (
	"errors"
	"net/http"

	"github.com/neo4j/cli/common/clictx"
	"github.com/neo4j/cli/neo4j/aura/internal/api"
	"github.com/neo4j/cli/neo4j/aura/internal/output"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	var (
		version              string
		region               string
		memory               string
		name                 string
		_type                string
		tenantId             string
		cloudProvider        string
		customerManagedKeyId string
	)

	const (
		versionFlag              = "version"
		regionFlag               = "region"
		memoryFlag               = "memory"
		nameFlag                 = "name"
		typeFlag                 = "type"
		tenantIdFlag             = "tenant-id"
		cloudProviderFlag        = "cloud-provider"
		customerManagedKeyIdFlag = "customer-managed-key-id"
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new instance",
		Long: `This subcommand starts the creation process of an Aura instance.

Creating an instance is an asynchronous operation. Supported instance configurations for your tenant can be obtained by calling the tenant get subcommand.

You can poll the current status of this operation by periodically getting the instance details for the instance ID using the get subcommand. Once the status transitions from "creating" to "running" you may begin to use your instance.

This subcommand returns your instance ID, initial credentials, connection URL along with your tenant id, cloud provider, region, instance type, and the instance name for you to use once the instance is running. It is important to store these initial credentials until you have the chance to login to your running instance and change them.

You must also provide a --cloud-provider flag with the subcommand, which specifies which cloud provider the instances will be hosted in. The acceptable values for this field are gcp, aws, or azure.

For Enterprise instances you can specify a --customer-managed-key-id flag to use a Customer Managed Key for encryption.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			typeValue, _ := cmd.Flags().GetString("type")
			if typeValue != "free-db" {
				cmd.MarkFlagRequired(memoryFlag)
			}

			config, ok := clictx.Config(cmd.Context())

			if !ok {
				return errors.New("error fetching cli configuration values")
			}

			if config.Aura.DefaultTenant == "" {
				cmd.MarkFlagRequired(tenantIdFlag)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{
				"version":        version,
				"region":         region,
				"name":           name,
				"type":           _type,
				"cloud_provider": cloudProvider,
			}

			if tenantId == "" {
				config, ok := clictx.Config(cmd.Context())
				if !ok {
					return errors.New("error fetching cli configuration values")
				}
				body["tenant_id"] = config.Aura.DefaultTenant
			} else {
				body["tenant_id"] = tenantId
			}

			if _type == "free-db" {
				body["memory"] = "1GB"
			} else {
				body["memory"] = memory
			}

			if customerManagedKeyId != "" {
				body["customer_managed_key_id"] = customerManagedKeyId
			}

			resBody, statusCode, err := api.MakeRequest(cmd, http.MethodPost, "/instances", body)
			if err != nil {
				return err
			}

			// NOTE: Instance create should not return OK (200), it always returns 202
			if statusCode == http.StatusAccepted || statusCode == http.StatusOK {
				err = output.PrintBody(cmd, resBody)
				if err != nil {
					return err
				}

			}

			return nil
		},
	}

	cmd.Flags().StringVar(&version, versionFlag, "5", "The Neo4j version of the instance.")

	cmd.Flags().StringVar(&region, regionFlag, "", "The region where the instance is hosted.")
	cmd.MarkFlagRequired(regionFlag)

	cmd.Flags().StringVar(&memory, memoryFlag, "", "The size of the instance memory in GB.")

	cmd.Flags().StringVar(&name, nameFlag, "", "The name of the instance (any UTF-8 characters with no trailing or leading whitespace).")
	cmd.MarkFlagRequired(nameFlag)

	cmd.Flags().StringVar(&_type, typeFlag, "", "The type of the instance.")
	cmd.MarkFlagRequired(typeFlag)

	cmd.Flags().StringVar(&tenantId, tenantIdFlag, "", "")

	cmd.Flags().StringVar(&cloudProvider, cloudProviderFlag, "", "The cloud provider hosting the instance.")
	cmd.MarkFlagRequired(cloudProviderFlag)

	cmd.Flags().StringVar(&customerManagedKeyId, customerManagedKeyIdFlag, "", "An optional customer managed key to be used for instance creation.")

	return cmd
}
