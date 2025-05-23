package clone

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	storageClassFlag      = "storage-class"
	storageSizeFlag       = "storage-size"
	recoveryTimestampFlag = "recovery-timestamp"
	recoveryDateFormat    = "2006-01-02T15:04:05-07:00"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId   string
	StorageClass *string
	StorageSize  *int64
	RecoveryDate *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("clone %s", instanceIdArg),
		Short: "Clones a PostgreSQL Flex instance",
		Long: "Clones a PostgreSQL Flex instance from a selected point in time. " +
			"The new cloned instance will be an independent instance with the same settings as the original instance unless the flags are specified.",
		Example: examples.Build(
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00`),
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage class.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-class premium-perf6-stackit`),
			examples.NewExample(
				`Clone a PostgreSQL Flex instance with ID "xxx" from a selected recovery timestamp and specify storage size.`,
				`$ stackit postgresflex instance clone xxx --recovery-timestamp 2023-04-17T09:28:00+00:00 --storage-size 10`),
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to clone instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("clone PostgreSQL Flex instance: %w", err)
			}
			instanceId := *resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Cloning instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.Region, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance cloning: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, instanceLabel, instanceId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(recoveryTimestampFlag, "", "Recovery timestamp for the instance, in a date-time with the layout format YYYY-MM-DDTHH:mm:ss±HH:mm, e.g. 2006-01-02T15:04:05-07:00")
	cmd.Flags().String(storageClassFlag, "", "Storage class. If not specified, storage class from the existing instance will be used.")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB). If not specified, storage size from the existing instance will be used.")

	err := flags.MarkFlagsRequired(cmd, recoveryTimestampFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	recoveryTimestamp, err := flags.FlagToDateTimePointer(p, cmd, recoveryTimestampFlag, recoveryDateFormat)
	if err != nil {
		return nil, &cliErr.FlagValidationError{
			Flag:    recoveryTimestampFlag,
			Details: err.Error(),
		}
	}
	recoveryTimestampString := recoveryTimestamp.Format(recoveryDateFormat)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		StorageClass:    flags.FlagToStringPointer(p, cmd, storageClassFlag),
		StorageSize:     flags.FlagToInt64Pointer(p, cmd, storageSizeFlag),
		RecoveryDate:    utils.Ptr(recoveryTimestampString),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

type PostgreSQLFlexClient interface {
	CloneInstance(ctx context.Context, projectId, region, instanceId string) postgresflex.ApiCloneInstanceRequest
	GetInstanceExecute(ctx context.Context, projectId, region, instanceId string) (*postgresflex.InstanceResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, region, flavorId string) (*postgresflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient PostgreSQLFlexClient) (postgresflex.ApiCloneInstanceRequest, error) {
	req := apiClient.CloneInstance(ctx, model.ProjectId, model.Region, model.InstanceId)

	var storages *postgresflex.ListStoragesResponse
	if model.StorageClass != nil || model.StorageSize != nil {
		currentInstance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.Region, model.InstanceId)
		if err != nil {
			return req, fmt.Errorf("get PostgreSQL Flex instance: %w", err)
		}
		validationFlavorId := currentInstance.Item.Flavor.Id
		currentInstanceStorageClass := currentInstance.Item.Storage.Class
		currentInstanceStorageSize := currentInstance.Item.Storage.Size

		storages, err = apiClient.ListStoragesExecute(ctx, model.ProjectId, model.Region, *validationFlavorId)
		if err != nil {
			return req, fmt.Errorf("get PostgreSQL Flex storages: %w", err)
		}

		if model.StorageClass == nil {
			err = postgresflexUtils.ValidateStorage(currentInstanceStorageClass, model.StorageSize, storages, *validationFlavorId)
		} else if model.StorageSize == nil {
			err = postgresflexUtils.ValidateStorage(model.StorageClass, currentInstanceStorageSize, storages, *validationFlavorId)
		} else {
			err = postgresflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, *validationFlavorId)
		}
		if err != nil {
			return req, err
		}
	}

	req = req.CloneInstancePayload(postgresflex.CloneInstancePayload{
		Class:     model.StorageClass,
		Size:      model.StorageSize,
		Timestamp: model.RecoveryDate,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, instanceLabel, instanceId string, resp *postgresflex.CloneInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("response not set")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgresFlex instance clone: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal PostgresFlex instance clone: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Cloned"
		if async {
			operationState = "Triggered cloning of"
		}
		p.Info("%s instance from instance %q. New Instance ID: %s\n", operationState, instanceLabel, instanceId)
		return nil
	}
}
