package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	userIdArg = "USER_ID"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
	UserId     string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", userIdArg),
		Short: "Shows details of a Secrets Manager user",
		Long:  "Shows details of a Secrets Manager user.",
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Secrets Manager user with ID "xxx" of instance with ID "yyy"`,
				"$ stackit secrets-manager user describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of a Secrets Manager user with ID "xxx" of instance with ID "yyy" in JSON format`,
				"$ stackit secrets-manager user describe xxx --instance-id yyy --output-format json"),
		),
		Args: args.SingleArg(userIdArg, utils.ValidateUUID),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Secrets Manager user: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, *resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	userId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		UserId:          userId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiGetUserRequest {
	req := apiClient.GetUser(ctx, model.ProjectId, model.InstanceId, model.UserId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, user secretsmanager.User) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(user, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager user: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(user.Id))
		table.AddSeparator()
		table.AddRow("USERNAME", utils.PtrString(user.Username))
		table.AddSeparator()
		if user.Description != nil && *user.Description != "" {
			table.AddRow("DESCRIPTION", *user.Description)
			table.AddSeparator()
		}
		if user.Password != nil && *user.Password != "" {
			table.AddRow("PASSWORD", *user.Password)
			table.AddSeparator()
		}
		table.AddRow("WRITE ACCESS", utils.PtrString(user.Write))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
