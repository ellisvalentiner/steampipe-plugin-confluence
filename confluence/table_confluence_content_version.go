package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentVersion() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_version",
		Description: "Confluence Content Version.",
		List: &plugin.ListConfig{
			ParentHydrate: listContentForVersions,
			Hydrate:       listContentVersion,
			KeyColumns:    plugin.OptionalColumns([]string{"id"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getContentVersion,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the content.",
			},
			{
				Name:        "number",
				Type:        proto.ColumnType_INT,
				Description: "The content version number.",
			},
			{
				Name:        "when",
				Type:        proto.ColumnType_STRING,
				Description: "When the content version was published.",
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "The content version message.",
			},
			{
				Name:        "minor_edit",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the version corresponds to a minor edit.",
			},
			{
				Name:        "username",
				Type:        proto.ColumnType_STRING,
				Description: "The username for the content version's author.",
			},
			{
				Name:        "userKey",
				Type:        proto.ColumnType_STRING,
				Description: "The user key for the content version's author.",
			},
			{
				Name:        "accountId",
				Type:        proto.ColumnType_STRING,
				Description: "The account ID for the content version's author.",
			},
			{
				Name:        "email",
				Type:        proto.ColumnType_STRING,
				Description: "The email for the content version's author.",
			},
			{
				Name:        "displayName",
				Type:        proto.ColumnType_STRING,
				Description: "The display name for the content version's author.",
			},
		},
	}
}

type contentVersion struct {
	ID          string `json:"id,omitempty"`
	Number      int    `json:"number,omitempty"`
	When        string `json:"when,omitempty"`
	Message     string `json:"message,omitempty"`
	MinorEdit   bool   `json:"minorEdit,omitempty"`
	Username    string `json:"username,omitempty"`
	UserKey     string `json:"userKey,omitempty"`
	AccountID   string `json:"accountId,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

//// LIST FUNCTIONS

func listContentVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentVersion")

	content := h.Item.(*model.ContentScheme)
	if content == nil || content.Version == nil || content.Version.By == nil {
		return nil, nil
	}
	row := contentVersion{
		ID:          content.ID,
		Number:      content.Version.Number,
		When:        content.Version.When,
		Message:     content.Version.Message,
		MinorEdit:   content.Version.MinorEdit,
		Username:    content.Version.By.Username,
		UserKey:     content.Version.By.UserKey,
		AccountID:   content.Version.By.AccountID,
		Email:       content.Version.By.Email,
		DisplayName: content.Version.By.DisplayName,
	}
	d.StreamListItem(ctx, row)

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getContentVersion(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getContentVersion")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	id := getStringQual(d, "id")
	logger.Trace("getContentVersion", "id", id)
	if id == "" {
		return nil, nil
	}

	expand := []string{}
	start := 0
	pageSize := 50
	var rows []contentVersion
	for {
		versions, _, err := instance.Content.Version.Gets(ctx, id, expand, start, pageSize)
		if err != nil {
			return nil, err
		}
		if versions == nil || len(versions.Results) == 0 {
			break
		}
		for _, version := range versions.Results {
			if plugin.IsCancelled(ctx) {
				return rows, nil
			}
			row := contentVersion{
				ID:        id,
				Number:    version.Number,
				When:      version.When,
				Message:   version.Message,
				MinorEdit: version.MinorEdit,
			}
			if version.By != nil {
				row.Username = version.By.Username
				row.UserKey = version.By.UserKey
				row.AccountID = version.By.AccountID
				row.Email = version.By.Email
				row.DisplayName = version.By.DisplayName
			}
			rows = append(rows, row)
		}
		if len(versions.Results) < pageSize {
			break
		}
		start += len(versions.Results)
	}

	return rows, nil
}

func listContentForVersions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return listContentWithExpand(ctx, d, []string{"version"})
}
