data "azuread_client_config" "current" {}
data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "dev" {
  location = local.location
  name     = "${local.project_name}-dev-rg"
}