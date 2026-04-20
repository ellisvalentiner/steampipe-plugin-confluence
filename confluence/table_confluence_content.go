package confluence

import (
	"context"
	"fmt"

	model "github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableConfluenceContent() *plugin.Table {
	return &plugin.Table{
		Name:        "confluence_content",
		Description: "Confluence Content.",
		List: &plugin.ListConfig{
			Hydrate: listContent,
			KeyColumns: plugin.OptionalColumns([]string{"id", "space_key", "type", "status", "title"}),
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
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: "The content title",
			},
			{
				Name:        "space_key",
				Type:        proto.ColumnType_STRING,
				Description: "The space containing the content",
				Transform:   transform.FromField("Space.Key"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The content status",
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The content type (page, blogpost, attachment or content)",
			},
			{
				Name:        "version_number",
				Type:        proto.ColumnType_INT,
				Description: "The content version",
				Transform:   transform.FromField("Version.Number"),
			},
			{
				Name:        "last_modified",
				Type:        proto.ColumnType_STRING,
				Description: "When the content was last modified.",
				Transform:   transform.FromField("Version.When"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listContent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return listContentWithExpand(ctx, d, []string{"space", "version"})
}

func listContentWithExpand(ctx context.Context, d *plugin.QueryData, expand []string) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("List confluence content")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	maxItems := getQueryLimit(d)
	contentID := getStringQual(d, "id")
	if contentID == "" {
		contentID = getStringQual(d, "content_id")
	}
	spaceKey := getStringQual(d, "space_key")
	contentType := getStringQual(d, "type")
	status := getStringQual(d, "status")
	title := getStringQual(d, "title")

	if contentID != "" {
		cql := fmt.Sprintf("id=%s", contentID)
		page, _, err := instance.Content.Search(ctx, cql, "", expand, "", 1)
		if err != nil {
			return nil, err
		}
		if page == nil || len(page.Results) == 0 {
			return nil, nil
		}

		d.StreamListItem(ctx, page.Results[0])
		return nil, nil
	}

	statuses := []string{}
	if status != "" {
		statuses = append(statuses, status)
	}

	options := &model.GetContentOptionsScheme{
		ContextType: contentType,
		SpaceKey:    spaceKey,
		Title:       title,
		Status:      statuses,
		Expand:      expand,
	}

	startAt := 0
	defaultPageSize := 100
	streamed := 0
	for {
		pageSize := requestPageSize(maxItems, streamed, defaultPageSize)
		if pageSize == 0 {
			break
		}

		page, response, err := instance.Content.Gets(ctx, options, startAt, pageSize)
		if err != nil {
			logger.Warn("Encountered error", "error", err, "Response", response)
			return nil, err
		}

		logger.Trace("Adding content items", "start", page.Start, "size", page.Size, "links", page.Links)
		for _, content := range page.Results {
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

func getContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("Get confluence content")

	instance, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	id := getStringQual(d, "id")
	logger.Trace("getContent", "id", id)
	if id == "" {
		return nil, nil
	}

	expand := []string{"space", "version", "body.storage", "body.view"}
	cql := fmt.Sprintf("id=%s", id)
	page, _, err := instance.Content.Search(ctx, cql, "", expand, "", 1)
	if err != nil {
		return nil, err
	}

	if page == nil || len(page.Results) == 0 {
		return nil, nil
	}

	return page.Results[0], nil
}
