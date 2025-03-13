terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.14.0"
    }
  }
}

resource "azurerm_storage_account" "storage" {
  account_replication_type = "LRS"
  account_tier             = "Standard"
  location                 = local.location
  name                     = "${local.project_shortname}${var.env_name}storagekjeld"
  resource_group_name      = data.azurerm_resource_group.rg.name
}

resource "azurerm_storage_container" "user-files" {
  name               = "user-files"
  storage_account_id = azurerm_storage_account.storage.id
}