---
organization: ellisvalentiner
category: ["software development"]
icon_url: "/images/plugins/ellisvalentiner/confluence.svg"
brand_color: "#2684FF"
display_name: "Confluence"
short_name: "confluence"
description: "Steampipe plugin for querying pages, spaces, and more from Confluence."
og_description: "Query Confluence with SQL! Open source CLI. No DB required."
og_image: ""
---

# Confluence + Steampipe

[Confluence](https://www.atlassian.com/software/confluence) is a collaboration wiki tool used to help teams to collaborate and share knowledge efficiently.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  id,
  title,
  space_key,
  status,
  type,
  version_number
from
  confluence_content;
```

```
+--------+------------------------------+-----------+---------+------+----------------+
| id     | title                        | space_key | status  | type | version_number |
+--------+------------------------------+-----------+---------+------+----------------+
| 163576 | Documentation                | DOC       | current | page | 5              |
| 110222 | HTTP Status Codes            | DOC       | current | page | 3              |
| 336504 | Confluence Getting started   | DOC       | current | page | 2              |
| 916343 | Staff Directory              | DOC       | current | page | 2              |
| 196895 | Product Requirements         | DOC       | current | page | 2              |
+--------+------------------------------+-----------+---------+------+----------------+
```

## Documentation

- **[Table definitions & examples →](https://hub.steampipe.io/plugins/ellisvalentiner/confluence/tables)**

## Get started

### Install

Download and install the latest Confluence plugin:

```bash
steampipe plugin install ellisvalentiner/confluence
```

### Credentials

| Item        | Description                                                                                                                            |
| :---------- | :------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | Confluence requires an [API token](https://id.atlassian.com/manage-profile/security/api-tokens), site base url and username for all requests. |

### Configuration

Installing the latest Confluence plugin will create a config file (`~/.steampipe/config/confluence.spc`) with a single connection named `confluence`:

```hcl
connection "confluence" {
    plugin    = "ellisvalentiner/confluence"

    # Base URI of your Confluence Cloud instance
    base_url = "https://your-domain.atlassian.net/"

    # The user name to access the Confluence cloud instance
    # username = "name@company.domain"

    # Access Token for the API
    # See https://id.atlassian.com/manage/api-tokens
    # token = ""
}
```

- `base_url` - The site url of your Atlassian subscription.
- `username` - Email address of agent user who have permission to access the API.
- `token` - [API token](https://id.atlassian.com/manage-profile/security/api-tokens) for user's Atlassian account.

## Get involved

- Open source: https://github.com/ellisvalentiner/steampipe-plugin-confluence
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
