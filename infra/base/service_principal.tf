resource "azuread_application" "this" {
  display_name = "${local.project_name} GitHub Actions SP"
  owners       = [data.azuread_client_config.current.object_id]
}

resource "azuread_service_principal" "this" {
  client_id                    = azuread_application.this.client_id
  app_role_assignment_required = false
  owners                       = [data.azuread_client_config.current.object_id]
}

resource "azurerm_role_assignment" "dev" {
  principal_id         = azuread_service_principal.this.object_id
  role_definition_name = "Contributor"
  scope                = azurerm_resource_group.dev.id
}

resource "time_rotating" "password_rotation" {
  rotation_days = 30
}

resource "azuread_service_principal_password" "this" {
  service_principal_id = azuread_service_principal.this.id
  rotate_when_changed = {
    rotation = time_rotating.password_rotation.id
  }
}