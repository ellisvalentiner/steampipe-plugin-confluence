package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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

	cql := getStringQual(d, "cql")
	if cql == "" {
		return nil, nil
	}
	logger.Trace("listSearchContent", "cql", cql)

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	maxItems := getQueryLimit(d)
	options := &model.SearchContentOptions{
		Expand: []string{"space"},
	}

	startAt := 0
	defaultPageSize := 100
	streamed := 0
	for {
		pageSize := requestPageSize(maxItems, streamed, defaultPageSize)
		if pageSize == 0 {
			break
		}

		options.Start = startAt
		options.Limit = pageSize

		searchResults, response, err := instance.Search.Content(ctx, cql, options)
		if err != nil {
			logger.Warn("Encountered error", "error", err, "Response", response)
			return nil, err
		}

		for _, row := range searchResults.Results {
			d.StreamListItem(ctx, row)
			streamed++
			if plugin.IsCancelled(ctx) {
				logger.Trace("CANCELLED!")
				return nil, nil
			}
			if maxItems > 0 && streamed >= maxItems {
				return nil, nil
			}
		}

		if searchResults.Size == 0 || searchResults.Size < options.Limit {
			break
		}

		startAt += searchResults.Size
	}
	return nil, nil
}
