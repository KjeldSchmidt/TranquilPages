resource "azurerm_key_vault" "this" {
  location            = local.location
  name                = "${local.project_name}-${var.env_name}-kv"
  resource_group_name = data.azurerm_resource_group.rg.name
  sku_name            = "standard"
  tenant_id           = data.azurerm_subscription.current.tenant_id
}

resource "azurerm_key_vault_secret" "database_connection_string" {
  key_vault_id = azurerm_key_vault.this.id
  name         = "database-connection-string"
  value        = azurerm_cosmosdb_account.this.primary_mongodb_connection_string
}