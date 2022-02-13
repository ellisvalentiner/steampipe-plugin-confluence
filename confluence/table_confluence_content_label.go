package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getContentLabel,
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

//// LIST FUNCTIONS

func listContentLabel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContentLabel")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	var maxResults int
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(100) {
			maxResults = int(*limit)
		}
	} else {
		maxResults = 100
	}

	startAt := 0
	content := h.Item.(*model.ContentScheme)
	quals := d.KeyColumnQuals
	prefix := quals["prefix"].GetStringValue()
	labels, _, err := instance.Content.Label.Gets(context.Background(), content.ID, prefix, startAt, maxResults)
	if err != nil {
		return nil, err
	}
	for _, label := range labels.Results {
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
		if plugin.IsCancelled(ctx) {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getContentLabel(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getContentLabel")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	var maxResults int
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(100) {
			maxResults = int(*limit)
		}
	} else {
		maxResults = 100
	}

	startAt := 0

	content := h.Item.(*model.ContentScheme)
	quals := d.KeyColumnQuals
	logger.Warn("getContentLabel", "quals", quals)
	contentID := quals["content_id"].GetStringValue()
	prefix := quals["prefix"].GetStringValue()

	labels, _, err := instance.Content.Label.Gets(context.Background(), contentID, prefix, startAt, maxResults)
	if err != nil {
		return nil, err
	}
	var rows []contentLabel

	for _, label := range labels.Results {
		row := contentLabel{
			ID:        label.ID,
			ContentID: content.ID,
			Title:     content.Title,
			SpaceKey:  content.Space.Key,
			Prefix:    label.Prefix,
			Name:      label.Name,
			Label:     label.Label,
		}
		// d.StreamListItem(ctx, row)
		if plugin.IsCancelled(ctx) {
			return nil, nil
		}
		rows = append(rows, row)
	}

	return rows, nil
}
