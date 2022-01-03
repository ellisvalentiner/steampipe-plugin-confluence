# Table: confluence_content

Content in a Confluence instance.

## Examples

### Get basic info about the content

```sql
select
  id,
  title,
  status,
  type
from
  confluence_content;
```

### Get the count of content type

```sql
select
  "type",
  count(*)
from
  confluence_content
group by "type";
```

### Get content with `draft` in the title

```sql
select
  *
from
  confluence_content
where
  title ilike '%draft%';
```
