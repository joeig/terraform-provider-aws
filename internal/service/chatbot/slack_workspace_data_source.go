// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package chatbot

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/chatbot"
	awstypes "github.com/aws/aws-sdk-go-v2/service/chatbot/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkDataSource("aws_chatbot_slack_workspace", name="Slack Workspace")
func newSlackWorkspaceDataSource(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &slackWorkspaceDataSource{}, nil
}

const (
	DSNameSlackWorkspace = "Slack Workspace Data Source"
)

type slackWorkspaceDataSource struct {
	framework.DataSourceWithModel[slackWorkspaceDataSourceModel]
}

func (d *slackWorkspaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slack_team_id": schema.StringAttribute{
				Computed: true,
			},
			"slack_team_name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (d *slackWorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Meta().ChatbotClient(ctx)

	var data slackWorkspaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := findSlackWorkspaceByName(ctx, conn, data.SlackTeamName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.Chatbot, create.ErrActionReading, DSNameSlackWorkspace, data.SlackTeamName.String(), err),
			err.Error(),
		)
		return
	}

	data.SlackTeamID = flex.StringToFramework(ctx, out.SlackTeamId)
	data.SlackTeamName = flex.StringToFramework(ctx, out.SlackTeamName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findSlackWorkspaceByName(ctx context.Context, conn *chatbot.Client, slack_team_name string) (*awstypes.SlackWorkspace, error) {
	input := &chatbot.DescribeSlackWorkspacesInput{
		MaxResults: aws.Int32(10),
	}

	for {
		output, err := conn.DescribeSlackWorkspaces(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, workspace := range output.SlackWorkspaces {
			if aws.ToString(workspace.SlackTeamName) == slack_team_name {
				return &workspace, nil
			}
		}

		if output.NextToken == nil {
			break
		}
		input.NextToken = output.NextToken
	}
	// If we are here, then we need to return an error that the slack workspace was not found.
	return nil, create.Error(names.Chatbot, "missing", DSNameSlackWorkspace, slack_team_name, nil)
}

type slackWorkspaceDataSourceModel struct {
	framework.WithRegionModel
	SlackTeamName types.String `tfsdk:"slack_team_name"`
	SlackTeamID   types.String `tfsdk:"slack_team_id"`
}
