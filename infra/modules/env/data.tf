data "terraform_remote_state" "base" {
  backend = "azurerm"

  config = {
    resource_group_name  = "tfstate"
    storage_account_name = "tfstatekjeldschmidt"
    container_name       = "tfstate"
    key                  = "tranquil-pages.base.tfstate"
  }
}

data "azurerm_resource_group" "rg" {
  name = data.terraform_remote_state.base.outputs["${var.env_name}_resource_group"].name
}