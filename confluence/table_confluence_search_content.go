package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableConfluenceSearchContent() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_search_content",
		Description: "Confluence Search Content.",
		List: &plugin.ListConfig{
			Hydrate:    listSearchContent,
			KeyColumns: plugin.SingleColumn("cql"),
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Automatically assigned when the content is created.",
				Transform:   transform.FromField("Content.ID"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The content title",
				Transform:   transform.FromField("Content.Title"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The content status",
				Transform:   transform.FromField("Content.Status"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The content type (page, blogpost, attachment or content)",
				Transform:   transform.FromField("Content.Type"),
			},
			{
				Name:        "last_modified",
				Type:        proto.ColumnType_STRING,
				Description: "When the content was last modified",
				Transform:   transform.FromField("LastModified"),
			},
			{
				Name:        "cql",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("cql"),
				Description: "The Confluence query langauge.",
			},
		},
	}
}

//// LIST FUNCTIONS

func listSearchContent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("Search confluence content")

	cql := d.KeyColumnQuals["cql"].GetStringValue()
	logger.Trace("listSearchContent", "cql", cql)

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	var limit int
	if d.QueryContext.Limit != nil {
		limit = int(limit)
	} else {
		limit = 100
	}
	options := &model.SearchContentOptions{
		Limit: limit,
		Expand: []string{"space"},
	}

	startAt := 0
	pageSize := 25
	pagesLeft := true
	for pagesLeft {
		searchResults, response, err := instance.Search.Content(context.Background(), cql, options)
		if err != nil {
			logger.Warn("Encountered error", "error", err, "Response", response)
			return nil, nil
		}

		for _, row := range searchResults.Results {
			d.StreamListItem(ctx, row)
			if plugin.IsCancelled(ctx) {
				logger.Trace("CANCELLED!")
				return nil, nil
			}
		}
		pagesLeft = false
		startAt += pageSize
	}
	return nil, nil
}
