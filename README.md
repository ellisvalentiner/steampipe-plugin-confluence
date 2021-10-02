# Atlassian Plugin for Steampipe

## Quick start

Install the plugin with Steampipe:

```sh
steampipe plugin install ellisvalentiner/atlassian
```

Run a query:

```sql
select
  id,
  title,
  status,
  type
from
  confluence_content
```

## Developing

Prerequisites:

* Steampipe
* Golang

Clone:

```sh
git clone https://github.com/ellisvalentiner/steampipe-plugin-atlassian.git
cd steampipe-plugin-atlassian
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
nano ~/.steampipe/config/atlassian.spc
```

Try it!

```sh
steampipe query
> .inspect atlassian
```
