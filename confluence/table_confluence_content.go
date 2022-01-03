package confluence

import (
	"context"

	"github.com/ctreminiom/go-atlassian/confluence"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableConfluenceContent() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content",
		Description: "Confluence Content.",
		List: &plugin.ListConfig{
			Hydrate: listContent,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getContent,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Automatically assigned when the content is created",
			},
			{
				Name:        "body",
				Type:        proto.ColumnType_JSON,
				Description: "The body of the content.",
			},
			{
				Name:        "child_types",
				Type:        proto.ColumnType_JSON,
				Description: "Shows whether a piece of content has attachments, comments, or child pages. Note, this doesn't actually contain the child objects.",
			},
			{
				Name:        "links",
				Type:        proto.ColumnType_JSON,
				Description: "",
			},
			{
				Name:        "space",
				Type:        proto.ColumnType_JSON,
				Description: "The space containing the content",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The content status",
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The content title",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The content type (page, blogpost, attachment or content)",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_JSON,
				Description: "The content version",
			},
		},
	}
}

//// LIST FUNCTIONS

func listContent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContent")

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
	// maxResults := 50

	options := &confluence.GetContentOptionsScheme{
		Expand: []string{"childTypes.all", "body.storage", "space", "version"},
	}

	pagesLeft := true
	for pagesLeft {
		page, _, err := instance.Content.Gets(context.Background(), options, startAt, maxResults)
		if err != nil {
			return nil, err
		}
		for _, content := range page.Results {
			d.StreamListItem(ctx, content)
			if plugin.IsCancelled(ctx) {
				return  nil, nil
			}
		}
		if page.Size < page.Limit {
			pagesLeft = false
		}
		startAt += maxResults
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getContent")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	logger.Warn("getContent", "quals", quals)
	id := quals["id"].GetStringValue()
	logger.Warn("getContent", "id", id)

	expand := []string{"childTypes.all", "body.storage", "space", "version"}
	version := 1

	content, _, err := instance.Content.Get(context.Background(), id, expand, version)
	if err != nil {
		return nil, err
	}

	return content, nil
}
