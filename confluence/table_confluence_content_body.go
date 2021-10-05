package confluence

import (
	"context"

	"github.com/ctreminiom/go-atlassian/confluence"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContentBody() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_body",
		Description: "Confluence Content Body.",
		List: &plugin.ListConfig{
			ParentHydrate: listContent,
			Hydrate:       listContentBody,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "",
			},
			{
				Name:        "view",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
			{
				Name:        "export_view",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
			{
				Name:        "styled_view",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
			{
				Name:        "storage",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
			{
				Name:        "editor2",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
			{
				Name:        "anonymous_export_view",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
		},
	}
}

// structs
type contentBody struct {
	ID                  string                     `json:"id,omitempty"`
	View                *confluence.BodyNodeScheme `json:"view"`
	ExportView          *confluence.BodyNodeScheme `json:"export_view"`
	StyledView          *confluence.BodyNodeScheme `json:"styled_view"`
	Storage             *confluence.BodyNodeScheme `json:"storage"`
	Editor2             *confluence.BodyNodeScheme `json:"editor2"`
	AnonymousExportView *confluence.BodyNodeScheme `json:"anonymous_export_view"`
}

//// LIST FUNCTIONS

func listContentBody(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentBody")

	content := h.Item.(*confluence.ContentScheme)
	c := contentBody{
		ID:                  content.ID,
		View:                content.Body.View,
		ExportView:          content.Body.ExportView,
		StyledView:          content.Body.StyledView,
		Storage:             content.Body.Storage,
		Editor2:             content.Body.Editor2,
		AnonymousExportView: content.Body.AnonymousExportView,
	}
	d.StreamListItem(ctx, c)

	return nil, nil
}
