# Table: confluence_search_content

Search content in a Confluence instance.

## Examples

### Get content with type "blogpost"

```sql
select
  id,
  title,
  status,
  type,
  last_modified
from
  confluence_search_content
where cql='type=blogpost';
```

### Get content with the "soc2" label

```sql
select
  id,
  title,
  status,
  type,
  last_modified
from
  confluence_search_content
where cql='label=soc2';
```
