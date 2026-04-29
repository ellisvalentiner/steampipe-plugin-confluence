# Table: confluence_content_version

Content versions in a Confluence instance.

## Examples

### Get basic info about the version

```sql
select
  *
from
  confluence_content_version;
```

### Get version count per content item

```sql
select
  id,
  count(*) as version_count
from
  confluence_content_version
group by id;
```

### Get the 50 oldest pages (join with confluence_content for title and space)

```sql
select
  c.title,
  c.space_key,
  v."when"
from
  confluence_content_version v
  join confluence_content c using (id)
order by
  v."when" asc
limit 50;
```
