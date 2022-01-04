# Table: confluence_content_body_storage

Content body from a Confluence instance in the `storage` format

## Examples

### Get the body for content in the storage format

```sql
select
  *
from
  confluence_content_body_storage;
```

### Get content body in the `storage` format

```sql
select
  *
from
  confluence_content_body_storage
where
  representation = 'storage';
```

### Get content that mentions `steampipe` in the body

```sql
select
  id,
  value
from
  confluence_content_body_storage
where
  "value" ilike '%steampipe%';
```
