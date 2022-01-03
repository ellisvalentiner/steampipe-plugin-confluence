# Table: confluence_content_body

Content body in a Confluence instance.

## Examples

### Get basic info about the content body

```sql
select
  id,
  storage
from
  confluence_content_body;
```

### Get content body in the `storage` format

```sql
select
  s.value
from
  confluence_content_body,
  jsonb_to_record(storage) AS s(value text);
```

### Get content body that mention `steampipe`

```sql
select
  s.value
from
  confluence_content_body,
  jsonb_to_record(storage) AS s(value text)
where
  s.value ilike '%steampipe%';
```
