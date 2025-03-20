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

resource "azurerm_key_vault_secret" "user_login_oauth_client_id" {
  key_vault_id = azurerm_key_vault.this.id
  name         = "user-login-oauth-client-id"
  value        = "Placeholder"

  lifecycle {
    ignore_changes = [value]
  }
}

resource "azurerm_key_vault_secret" "user_login_oauth_client_secret" {
  key_vault_id = azurerm_key_vault.this.id
  name         = "user-login-oauth-client-secret"
  value        = "Placeholder"

  lifecycle {
    ignore_changes = [value]
  }
}

resource "azurerm_key_vault_access_policy" "pipeline_service_principal" {
  key_vault_id = azurerm_key_vault.this.id
  object_id    = data.terraform_remote_state.base.outputs["pipeline_service_principal"].object_id
  tenant_id    = data.azurerm_subscription.current.tenant_id

  secret_permissions = ["Get", "List", "Set", "Delete"]
}