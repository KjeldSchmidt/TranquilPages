resource "azurerm_container_app_environment" "this" {
  name                = "${local.project_name}-${var.env_name}-container-app-env"
  resource_group_name = data.azurerm_resource_group.rg.name
  location            = data.azurerm_resource_group.rg.location
}

resource "azurerm_container_app" "this" {
  name                         = "${local.project_name}-${var.env_name}-container-app"
  resource_group_name          = data.azurerm_resource_group.rg.name
  container_app_environment_id = azurerm_container_app_environment.this.id
  revision_mode                = "Single"

  template {
    min_replicas = 0
    max_replicas = 1

    container {
      name   = "${local.project_name}-${var.env_name}-container"
      image  = "docker.io/kjeldschmidt2/tranquil-pages:latest"
      cpu    = "0.5"
      memory = "0.5Gi"

      env {
        name        = "DB_URL"
        secret_name = azurerm_key_vault_secret.database_connection_string.name
      }
    }
  }

  secret {
    name                = azurerm_key_vault_secret.database_connection_string.name
    identity            = "System"
    key_vault_secret_id = azurerm_key_vault_secret.database_connection_string.key_vault_id
  }
}