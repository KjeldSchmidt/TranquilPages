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

  identity {
    type = "SystemAssigned"
  }

  template {
    min_replicas = 0
    max_replicas = 1

    container {
      name   = "${local.project_name}-${var.env_name}-container"
      image  = "docker.io/kjeldschmidt2/tranquil-pages:latest"
      cpu    = "0.25"
      memory = "0.5Gi"

      env {
        name        = "DB_URL"
        secret_name = azurerm_key_vault_secret.database_connection_string.name
      }

      env {
        name        = "OAUTH_CLIENT_ID"
        secret_name = azurerm_key_vault_secret.user_login_oauth_client_id.name
      }

      env {
        name        = "OAUTH_CLIENT_SECRET"
        secret_name = azurerm_key_vault_secret.user_login_oauth_client_secret.name
      }

      env {
        name        = "JWT_SECRET"
        secret_name = azurerm_key_vault_secret.jwt_secret.name
      }

      env {
        name  = "FRONTEND_URL"
        value = azurerm_storage_account.storage.primary_web_endpoint
      }
    }
  }

  ingress {
    target_port      = 8080
    external_enabled = true

    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }

  secret {
    name                = azurerm_key_vault_secret.database_connection_string.name
    key_vault_secret_id = azurerm_key_vault_secret.database_connection_string.id
    identity            = "System"
  }

  secret {
    name                = azurerm_key_vault_secret.user_login_oauth_client_id.name
    key_vault_secret_id = azurerm_key_vault_secret.user_login_oauth_client_id.versionless_id
    identity            = "System"
  }

  secret {
    name                = azurerm_key_vault_secret.user_login_oauth_client_secret.name
    key_vault_secret_id = azurerm_key_vault_secret.user_login_oauth_client_secret.versionless_id
    identity            = "System"
  }

  secret {
    name                = azurerm_key_vault_secret.jwt_secret.name
    key_vault_secret_id = azurerm_key_vault_secret.jwt_secret.versionless_id
    identity            = "System"
  }
}

resource "azurerm_key_vault_access_policy" "container_app" {
  key_vault_id = azurerm_key_vault.this.id
  tenant_id    = data.azurerm_subscription.current.tenant_id
  object_id    = azurerm_container_app.this.identity[0].principal_id

  secret_permissions = ["Get"]
}
