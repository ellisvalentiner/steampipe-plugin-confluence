# Table: confluence_content_label

Content labels in a Confluence instance.

## Examples

### Get basic info about the labels

```sql
select
  id,
  content_id,
  title,
  space_key,
  prefix,
  name,
  label
from
  confluence_content_label;
```

### Get the count of content by label

```sql
select
  label,
  count(*)
from
  confluence_content_label
group by label;
```

### Get labels with `documentation` in the name

```sql
select
  *
from
  confluence_content_label
where
  name ilike '%documentation%';
```
