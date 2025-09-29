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
  instance_type    = "c1.small"    # Instance type
  location         = "NYC1"         # Location code
  operating_system = "Ubuntu 24.04" # Operating system
  hostname         = "my-server"
  friendly_name    = "My Test Server"
}
```

The provider automatically:
- Validates that the instance type exists
- Checks inventory availability in the specified location
- Confirms the operating system is available for that instance type and location
- Provides helpful error messages with alternatives if anything is invalid


## Contributing

Interested in contributing? See our [Contributing Guide](CONTRIBUTING.md) for development setup, testing, and guidelines.

## Example

See the `examples/` directory for complete usage examples.

```bash
cd examples/
terraform init
terraform plan
terraform apply
```

Remember to set your `ICS_API_TOKEN` environment variable before running Terraform commands.