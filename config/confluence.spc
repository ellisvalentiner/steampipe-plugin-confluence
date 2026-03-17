connection "confluence" {
    plugin = "ellisvalentiner/confluence"

    # --- Confluence Cloud (default) ---
    # deployment_type = "cloud"   # default; can be omitted
    # Base URI of your Confluence Cloud instance
    # base_url = "https://your-domain.atlassian.net/"
    # The user name to access the Confluence cloud instance
    # username = "name@company.domain"
    # token    = "your-cloud-api-token" # Atlassian Cloud API token

    # --- Confluence Data Center / Server ---
    # deployment_type = "datacenter"
    # Base URI of your Confluence Datacenter instance
    # base_url = "https://confluence.your-company.com/"
    # token    = "your-data-center-PAT" # Personal Access Token; username not required
}
