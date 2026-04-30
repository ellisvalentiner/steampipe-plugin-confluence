connection "confluence" {
    plugin = "ellisvalentiner/confluence"

    # --- Confluence Cloud (default) ---
    # deployment_type = "cloud"   # default; can be omitted
    # Base URI of your Confluence Cloud instance
    # Can also be set via the CONFLUENCE_BASE_URL environment variable.
    # base_url = "https://your-domain.atlassian.net/"
    # The user name (email) to access the Confluence cloud instance
    # Can also be set via the CONFLUENCE_USERNAME environment variable.
    # username = "name@company.domain"
    # Atlassian Cloud API token
    # Can also be set via the CONFLUENCE_TOKEN environment variable.
    # token    = "your-cloud-api-token"

    # --- Confluence Data Center / Server ---
    # deployment_type = "datacenter"
    # Can also be set via the CONFLUENCE_DEPLOYMENT_TYPE environment variable.
    # Base URI of your Confluence Datacenter instance
    # Can also be set via the CONFLUENCE_BASE_URL environment variable.
    # base_url = "https://confluence.your-company.com/"
    # Personal Access Token; username not required for Data Center
    # Can also be set via the CONFLUENCE_TOKEN environment variable.
    # token    = "your-data-center-PAT"
}
