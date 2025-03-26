resource "azurerm_storage_account_static_website" "frontend" {
  storage_account_id = azurerm_storage_account.storage.id
  index_document     = "index.html"
}