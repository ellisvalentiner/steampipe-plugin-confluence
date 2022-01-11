# Table: confluence_space

Spaces in a Confluence instance.

## Examples

### Get basic info about the spaces

```sql
select
  id,
  key,
  name,
  type,
  status
from
  confluence_space;
```

### Get the `global` spaces

```sql
select
  *
from
  confluence_space
where
  "type"='global';
```

### Get personal, archived spaces

```sql
select
  *
from
  confluence_space
where
  "type"='personal'
  and status='archived';
```
