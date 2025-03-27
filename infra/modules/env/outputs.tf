output "frontend_url" {
  value = azurerm_storage_account.storage.primary_web_endpoint
}

output "backend_url" {
  value = azurerm_container_app_environment.this.default_domain
}