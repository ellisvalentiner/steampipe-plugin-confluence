package confluence

import (
	"context"

	// "github.com/ctreminiom/go-atlassian/confluence"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	// "github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableConfluenceContentLabel() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content_label",
		Description: "Confluence Content Label.",
		// List: &plugin.ListConfig{
		// 	Hydrate: listContentLabel,
		// },
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

//// LIST FUNCTIONS

// func listContentLabel(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	logger.Trace("listContentLabel")
//
// 	instance, err := connect(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var maxResults int
// 	limit := d.QueryContext.Limit
// 	if limit != nil {
// 		if *limit < int64(100) {
// 			maxResults = int(*limit)
// 		}
// 	} else {
// 		maxResults = 100
// 	}
//
// 	startAt := 0
//
// 	pagesLeft := true
// 	for pagesLeft {
// 		labels, _, err := instance.Content.Label.Gets(context.Background(), contentID, prefix, startAt, maxResults)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, label := range labels.Results {
// 			d.StreamListItem(ctx, label)
// 			if plugin.IsCancelled(ctx) {
// 				return nil, nil
// 			}
// 		}
// 		if labels.Size < labels.Limit {
// 			pagesLeft = false
// 		}
// 		startAt += maxResults
// 	}
// 	return nil, nil
// }

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

	quals := d.KeyColumnQuals
	logger.Warn("getContentLabel", "quals", quals)
	// id := quals["id"].GetStringValue()
	// logger.Warn("getContentLabel", "id", id)
	contentID := quals["contentid"].GetStringValue()
	prefix := quals["prefix"].GetStringValue()

	content, _, err := instance.Content.Label.Gets(context.Background(), contentID, prefix, startAt, maxResults)
	if err != nil {
		return nil, err
	}

	return content, nil
}
