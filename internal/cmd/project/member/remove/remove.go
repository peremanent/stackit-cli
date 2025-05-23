package remove

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

const (
	roleFlag  = "role"
	forceFlag = "force"

	subjectArg = "SUBJECT"

	projectResourceType = "project"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Subject string
	Role    *string
	Force   bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("remove %s", subjectArg),
		Short: "Removes a member from a project",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Removes a member from a project.",
			"A member is a combination of a subject (user, service account or client) and a role.",
			"The subject is usually email address for users or name in case of clients",
		),
		Args: args.SingleArg(subjectArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Remove a member (user "someone@domain.com" with an "editor" role) from a project`,
				"$ stackit project member remove someone@domain.com --project-id xxx --role editor"),
			examples.NewExample(
				`Remove a member (user "someone@domain.com" with a "reader" role) from a project, along with all other roles of the subject that would stop the removal of the "reader" role`,
				"$ stackit project member remove someone@domain.com --project-id xxx --role reader --force"),
		),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to remove the role %q from %s on project %q?", *model.Role, model.Subject, projectLabel)
				if model.Force {
					prompt = fmt.Sprintf("%s This will also remove other roles of the subject that would stop the removal of the requested role", prompt)
				}
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("remove member: %w", err)
			}

			params.Printer.Info("Removed the role %q from %s on project %q\n", utils.PtrString(model.Role), model.Subject, projectLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(roleFlag, "", "The role to be removed from the subject")
	cmd.Flags().Bool(forceFlag, false, "When true, removes other roles of the subject that would stop the removal of the requested role")

	err := flags.MarkFlagsRequired(cmd, roleFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	subject := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Subject:         subject,
		Role:            flags.FlagToStringPointer(p, cmd, roleFlag),
		Force:           flags.FlagToBoolValue(p, cmd, forceFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiRemoveMembersRequest {
	req := apiClient.RemoveMembers(ctx, model.GlobalFlagModel.ProjectId)
	payload := authorization.RemoveMembersPayload{
		Members: utils.Ptr([]authorization.Member{
			{
				Subject: utils.Ptr(model.Subject),
				Role:    model.Role,
			},
		}),
		ResourceType: utils.Ptr(projectResourceType),
	}
	payload.ForceRemove = &model.Force
	req = req.RemoveMembersPayload(payload)
	return req
}
