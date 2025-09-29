# Terraform Provider for Ingenuity Cloud Services

This Terraform provider allows you to manage bare metal servers on Ingenuity Cloud Services (ICS) infrastructure.

## Features

- **Ultra-Simple Interface**: Just specify `instance_type`, `location`, and `operating_system` - the provider handles everything else
- **Automatic Validation**: Real-time validation with helpful error messages showing available alternatives
- **Zero Discovery Required**: No need for data sources or complex lookups - just specify what you want
- **Comprehensive Error Messages**: Clear guidance when combinations aren't available
- **Automatic Provisioning**: Polls for up to 30 minutes waiting for server provisioning to complete
- **Simplified Billing**: Automatically uses hourly billing for easy cleanup and testing
- **Optional Data Sources**: Available for discovery if needed, but not required for basic usage

## Installation

### Using Terraform Registry

```hcl
terraform {
  required_providers {
    ics = {
      source = "UK2Group/ics"
    }
  }
}
```

### Local Development

1. Clone this repository
2. Build the provider: `go build -o terraform-provider-ics`
3. Move the binary to your Terraform plugins directory

## Authentication

Set your ICS API token as an environment variable:

```bash
export ICS_API_TOKEN="your-api-token-here"
```

Or configure it directly in your Terraform configuration:

```hcl
provider "ics" {
  api_token = "your-api-token-here"
}
```

## Usage

### Data Sources

#### `ics_inventory`

Retrieves available server inventory with friendly names and location information.

```hcl
data "ics_inventory" "available_servers" {}

# Show only servers with auto-provision inventory
output "available_servers" {
  value = {
    for item in data.ics_inventory.available_servers.items :
    "${item.sku_product_name}-${item.location_code}" => {
      server_type    = item.sku_product_name
      location       = item.location_code
      auto_provision = item.auto_provision_quantity
      cpu_info       = "${item.cpu_brand} ${item.cpu_model}"
      ram_gb         = item.total_ram_gb
      price_hourly   = item.price_hourly
    }
    if item.auto_provision_quantity > 0
  }
}
```

#### `ics_operating_systems`

Retrieves available operating systems for a specific server type and location.

```hcl
data "ics_operating_systems" "example" {
  server_type_name = "c1i.small"
  location         = "NYC1"
}

output "available_os" {
  value = data.ics_operating_systems.example.operating_systems
}
```

### Resources

#### `ics_bare_metal_server`

Provisions a bare metal server with automatic validation.

```hcl
resource "ics_bare_metal_server" "example" {
  instance_type    = "c1i.small"    # Instance type
  location         = "NYC1"         # Location code
  operating_system = "Ubuntu 24.04" # Operating system
  hostname         = "my-server"
  domain           = "example.com"
  friendly_name    = "My Test Server"
  notes           = "Provisioned via Terraform"
}
```

The provider automatically:
- Validates that the instance type exists
- Checks inventory availability in the specified location
- Confirms the operating system is available for that instance type and location
- Provides helpful error messages with alternatives if anything is invalid

## API Workflow

The provider implements the following workflow for server provisioning:

1. **Resolve**: Look up SKU ID from friendly server type name and location code
2. **Validate**: Check auto-provision inventory availability for the location
3. **Discover**: Retrieve available operating systems for the server type and location
4. **Order**: POST to `/rest-api/server-orders/order` with resolved configuration and OS product code
5. **Poll**: GET `/rest-api/servers` repeatedly (every 30 seconds for up to 30 minutes)
6. **Complete**: When the server appears in the servers list, provisioning is complete

For server destruction:
- **Cancel**: DELETE `/rest-api/servers/{server_id}/cancel` (automatically supported since we use hourly billing)

## Key Improvements

- **User-Friendly**: No need to look up SKU IDs manually - just use friendly names like `c1i.small`
- **Location-Aware**: Automatically validates inventory availability in specified locations
- **OS Discovery**: Dynamically discovers available operating systems instead of hardcoding
- **Simplified Billing**: Always uses hourly billing for consistent, automated cleanup
- **Auto-Validation**: Ensures servers can only be provisioned where inventory is available

## Development

### Building

```bash
go build -o terraform-provider-ics
```

### Testing

```bash
go test ./...
```

## Example

See the `examples/` directory for complete usage examples.

```bash
cd examples/
terraform init
terraform plan
terraform apply
```

Remember to set your `ICS_API_TOKEN` environment variable before running Terraform commands.