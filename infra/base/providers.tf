provider "azurerm" {
  subscription_id = "61c1fbf4-07d4-48c7-9d95-81aff1db63a8"
  features {}
}

terraform {
  backend "azurerm" {
    resource_group_name  = "tfstate"
    storage_account_name = "tfstatekjeldschmidt"
    container_name       = "tfstate"
    key                  = "tranquil-pages.base.tfstate"
  }
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "3.0.2"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.12.1"
    }
  }
}