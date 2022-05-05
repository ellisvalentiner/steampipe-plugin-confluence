package confluence

import (
	"context"

	model "github.com/ctreminiom/go-atlassian/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableConfluenceSpace() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_space",
		Description: "Confluence Space.",
		List: &plugin.ListConfig{
			Hydrate: listSpace,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSpace,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Automatically assigned when the space is created",
			},
			{
				Name:        "key",
				Type:        proto.ColumnType_STRING,
				Description: "The key of the space.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the space.",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of space.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the space.",
			},
		},
	}
}

//// LIST FUNCTIONS

func listSpace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listContent")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	startAt := 0
	pageSize := 25

	quals := d.KeyColumnQuals
	options := &model.GetSpacesOptionScheme{
		SpaceKeys: nil,
		// Type:      quals["type"].GetStringValue(),
		Status: quals["status"].GetStringValue(),
	}

	pagesLeft := true
	for pagesLeft {
		page, _, err := instance.Space.Gets(context.Background(), options, startAt, pageSize)
		if err != nil {
			return nil, err
		}
		for _, content := range page.Results {
			d.StreamListItem(ctx, content)
			if plugin.IsCancelled(ctx) {
				return nil, nil
			}
		}
		if page.Size < page.Limit {
			pagesLeft = false
		}
		startAt += pageSize
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getSpace(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("getSpace")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	quals := d.KeyColumnQuals
	logger.Warn("getSpace", "quals", quals)
	id := quals["id"].GetStringValue()
	logger.Warn("getSpace", "id", id)

	content, _, err := instance.Space.Get(context.Background(), id, []string{})
	if err != nil {
		return nil, err
	}

	return content, nil
}
