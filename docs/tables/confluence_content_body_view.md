# Table: confluence_content_body_view

Content body from a Confluence instance in the `view` format

## Examples

### Get the body for content in the view format

```sql
select
  *
from
  confluence_content_body_view;
```

### Get content body in the `view` format

```sql
select
  *
from
  confluence_content_body_view
where
  representation = 'view';
```

### Get content that mentions `steampipe` in the body

```sql
select
  id,
  value
from
  confluence_content_body_view
where
  "value" ilike '%steampipe%';
```
