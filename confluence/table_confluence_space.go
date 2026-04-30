package confluence

import (
	"context"
	"fmt"

	model "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableConfluenceSpace() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_space",
		Description: "Confluence Space.",
		List: &plugin.ListConfig{
			Hydrate:    listSpace,
			KeyColumns: plugin.OptionalColumns([]string{"key", "type", "status"}),
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
	logger.Trace("listSpace")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	startAt := 0
	defaultPageSize := 100
	maxItems := getQueryLimit(d)
	streamed := 0

	spaceKey := getStringQual(d, "key")
	spaceType := getStringQual(d, "type")
	status := getStringQual(d, "status")

	spaceKeys := []string{}
	if spaceKey != "" {
		spaceKeys = append(spaceKeys, spaceKey)
	}

	options := &model.GetSpacesOptionScheme{
		SpaceKeys: spaceKeys,
		SpaceType: spaceType,
		Status:    status,
	}

	for {
		pageSize := requestPageSize(maxItems, streamed, defaultPageSize)
		if pageSize == 0 {
			break
		}

		page, _, err := instance.Space.Gets(ctx, options, startAt, pageSize)
		if err != nil {
			return nil, err
		}
		if page == nil {
			return nil, fmt.Errorf("confluence space list returned nil page")
		}
		if page.Results == nil {
			break
		}

		for _, content := range page.Results {
			if content == nil {
				continue
			}
			d.StreamListItem(ctx, content)
			streamed++
			if plugin.IsCancelled(ctx) {
				return nil, nil
			}
			if maxItems > 0 && streamed >= maxItems {
				return nil, nil
			}
		}

		if page.Size == 0 || page.Size < page.Limit || page.Size < pageSize {
			break
		}

		startAt += page.Size
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

	id := getStringQual(d, "id")
	logger.Trace("getSpace", "id", id)
	if id == "" {
		return nil, nil
	}

	content, _, err := instance.Space.Get(ctx, id, []string{})
	if err != nil {
		return nil, err
	}

	return content, nil
}
