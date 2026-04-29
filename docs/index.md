---
organization: ellisvalentiner
category: ["software development"]
icon_url: "/images/plugins/ellisvalentiner/confluence.svg"
brand_color: "#2684FF"
display_name: "Confluence"
short_name: "confluence"
description: "Steampipe plugin for querying pages, spaces, and more from Confluence."
og_description: "Query Confluence with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/ellisvalentiner/confluence-social-graphic.png"
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
  version_number,
  last_modified,
  tags
from
  confluence_content;
```

```
+--------+------------------------------+-----------+---------+------+----------------+----------------------+-------------------------------+
| id     | title                        | space_key | status  | type | version_number | last_modified        | tags                          |
+--------+------------------------------+-----------+---------+------+----------------+----------------------+-------------------------------+
| 163576 | Documentation                | DOC       | current | page | 5              | 2024-01-10T16:23:58Z | ["docs", "how-to"]          |
| 110222 | HTTP Status Codes            | DOC       | current | page | 3              | 2023-11-02T09:14:21Z | ["reference", "http"]       |
| 336504 | Confluence Getting started   | DOC       | current | page | 2              | 2023-09-18T20:41:07Z | ["getting-started"]          |
| 916343 | Staff Directory              | DOC       | current | page | 2              | 2023-08-07T13:05:44Z | ["people", "directory"]     |
| 196895 | Product Requirements         | DOC       | current | page | 2              | 2023-07-29T17:52:11Z | ["product", "requirements"] |
+--------+------------------------------+-----------+---------+------+----------------+----------------------+-------------------------------+
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
    # Can also be set via the CONFLUENCE_BASE_URL environment variable.
    base_url = "https://your-domain.atlassian.net/"

    # The user name (email) to access the Confluence cloud instance
    # Can also be set via the CONFLUENCE_USERNAME environment variable.
    # username = "name@company.domain"

    # Atlassian Cloud API token
    # See https://id.atlassian.com/manage/api-tokens
    # Can also be set via the CONFLUENCE_TOKEN environment variable.
    # token = ""
}
```

- `base_url` - The site url of your Atlassian subscription. Can also be set via the `CONFLUENCE_BASE_URL` environment variable.
- `username` - Email address of agent user who have permission to access the API. Can also be set via the `CONFLUENCE_USERNAME` environment variable.
- `token` - [API token](https://id.atlassian.com/manage-profile/security/api-tokens) for user's Atlassian account. Can also be set via the `CONFLUENCE_TOKEN` environment variable.
- `deployment_type` - Either `cloud` (default) or `datacenter`. Can also be set via the `CONFLUENCE_DEPLOYMENT_TYPE` environment variable.

## Get involved

- Open source: https://github.com/ellisvalentiner/steampipe-plugin-confluence
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
