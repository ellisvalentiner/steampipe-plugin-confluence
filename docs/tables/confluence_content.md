# Table: confluence_content

Content in a Confluence instance.

## Examples

### Get basic info about the content

```sql
select
  id,
  title,
  space_key,
  status,
  type,
  last_modified
from
  confluence_content;
```

### Get recently modified content

```sql
select
  id,
  title,
  space_key,
  last_modified
from
  confluence_content
where
  last_modified is not null
order by
  last_modified desc;
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
