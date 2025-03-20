resource "azurerm_cosmosdb_account" "this" {
  name                 = "${local.project_name}-${var.env_name}-cosmos"
  location             = data.azurerm_resource_group.rg.location
  resource_group_name  = data.azurerm_resource_group.rg.name
  kind                 = "MongoDB"
  offer_type           = "Standard"
  mongo_server_version = "7.0"

  consistency_policy {
    consistency_level = "Eventual"
  }

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "EnableServerless"
  }

  # Disable backup to reduce costs
  backup {
    type = "Continuous"
    tier = "Continuous7Days"
  }

  # Enable network access from Azure services and your IP
  network_acl_bypass_for_azure_services = true
  network_acl_bypass_ids                = [] # Add your IP addresses here if needed

  # Set minimum RU/s to 0 for serverless
  capacity {
    total_throughput_limit = 0
  }

  # Enable automatic failover for high availability
  geo_location {
    location          = data.azurerm_resource_group.rg.location
    failover_priority = 0
    zone_redundant    = false
  }

  # Enable public network access (can be restricted later)
  public_network_access_enabled = true
}

# Create the MongoDB database
resource "azurerm_cosmosdb_mongo_database" "this" {
  name                = "tranquil_pages"
  resource_group_name = data.azurerm_resource_group.rg.name
  account_name        = azurerm_cosmosdb_account.this.name
}

# Create the books collection
resource "azurerm_cosmosdb_mongo_collection" "books" {
  name                = "books"
  resource_group_name = data.azurerm_resource_group.rg.name
  account_name        = azurerm_cosmosdb_account.this.name
  database_name       = azurerm_cosmosdb_mongo_database.this.name

  # Define the collection schema
  index {
    keys   = ["_id"]
    unique = true
  }

  # Add indexes for common queries
  index {
    keys = ["title"]
  }
  index {
    keys = ["author"]
  }
}

# Create the oauth_states collection with TTL index
resource "azurerm_cosmosdb_mongo_collection" "oauth_states" {
  name                = "oauth_states"
  resource_group_name = data.azurerm_resource_group.rg.name
  account_name        = azurerm_cosmosdb_account.this.name
  database_name       = azurerm_cosmosdb_mongo_database.this.name
  default_ttl_seconds = 60 * 30

  # Define the collection schema
  index {
    keys   = ["_id"]
    unique = true
  }
}

resource "azurerm_cosmosdb_mongo_collection" "blacklisted_jwts" {
  name                = "blacklisted_jwts"
  resource_group_name = data.azurerm_resource_group.rg.name
  account_name        = azurerm_cosmosdb_account.this.name
  database_name       = azurerm_cosmosdb_mongo_database.this.name
  default_ttl_seconds = 60 * 60 * 24 * 180 # Keep the blacklist for a long time, in case our tokens ever get a much longer TTL

  # Define the collection schema
  index {
    keys   = ["_id"]
    unique = true
  }

  # Add index for token lookups
  index {
    keys = ["token"]
  }
}