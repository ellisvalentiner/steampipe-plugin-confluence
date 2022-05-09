package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentLabel() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_label",
		Description: "Confluence Content Label.",
		List: &plugin.ListConfig{
			ParentHydrate: listContent,
			Hydrate:       listContentLabel,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Automatically assigned when the label is created",
			},
			{
				Name:        "content_id",
				Type:        proto.ColumnType_STRING,
				Description: "Automatically assigned when the content is created",
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The content title",
			},
			{
				Name:        "space_key",
				Type:        proto.ColumnType_STRING,
				Description: "The space containing the content",
			},
			{
				Name:        "prefix",
				Type:        proto.ColumnType_STRING,
				Description: "The label prefix",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The label name",
			},
			{
				Name:        "label",
				Type:        proto.ColumnType_STRING,
				Description: "The label",
			},
		},
	}
}

type contentLabel struct {
	ID        string `json:"id,omitempty"`
	ContentID string `json:"contentId,omitempty"`
	Title     string `json:"title,omitempty"`
	SpaceKey  string `json:"spaceKey,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Name      string `json:"name,omitempty"`
	Label     string `json:"label,omitempty"`
}

// LIST FUNCTIONS
func listContentLabel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentLabel")

	content := h.Item.(*model.ContentScheme)
	for _, label := range content.Metadata.Labels.Results {
		row := contentLabel{
			ID:        label.ID,
			ContentID: content.ID,
			Title:     content.Title,
			SpaceKey:  content.Space.Key,
			Prefix:    label.Prefix,
			Name:      label.Name,
			Label:     label.Label,
		}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}
