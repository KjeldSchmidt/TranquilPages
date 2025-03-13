output "client_id" {
  value = azuread_service_principal.this.id
}

output "client_secret" {
  value     = azuread_service_principal_password.this.value
  sensitive = true
}

output "subscription_id" {
  value = data.azurerm_subscription.current.subscription_id
}

output "tenant_id" {
  value = data.azurerm_subscription.current.tenant_id
}

output "dev_resource_group" {
  value = azurerm_resource_group.dev
}