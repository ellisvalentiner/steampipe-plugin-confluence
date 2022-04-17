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

### Get the count of content by label

```sql
select
  label,
  count(*)
from
  confluence_content_version
group by label;
```

### Get the 50 oldest pages

```sql
select
  title,
  space_key,
  "when"
from confluence_content_version
join confluence_content using (id)
order by "when" asc
limit 50;
```
