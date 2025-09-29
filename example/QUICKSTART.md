# Quick Start Guide

This is the fastest way to get the ICS Terraform provider working with a full example.

## 1. Build and Install the Provider

From the root directory of the repository:

```bash
make install
```

## 2. Set Your API Token

```bash
export ICS_API_TOKEN="your-ics-api-token-here"
```

## 3. Run the Example

```bash
cd full-example
./run-example.sh
```

## What the Example Does

1. **Validates your configuration** - Automatically checks that your instance type, location, and OS are available
2. **Provisions a server** - Orders and waits for the server to be provisioned (5-30 minutes)
3. **Shows results** - Displays the service ID immediately and public IP after provisioning

The provider handles all validation in the background with helpful error messages if something isn't available.

## Key Outputs

- **Service ID**: Available immediately after ordering (usually within seconds)
- **Public IP**: Available after the server is fully provisioned (5-30 minutes)
- **Server Details**: Complete information about the provisioned server

## Example Configuration

The example uses this simple configuration:

```hcl
resource "ics_bare_metal_server" "example" {
  instance_type    = "c1.small"    # Instance type
  location         = "NYC1"         # Location code
  operating_system = "Ubuntu 24.04" # Operating system
  hostname         = "ics-terraform-demo"
  domain           = "example.com"
  friendly_name    = "ICS Terraform Demo Server"
  notes           = "Demo server provisioned via Terraform"
}
```

The provider automatically validates all values and provides helpful error messages if any combination isn't available.

## Cleanup

⚠️ **Don't forget to destroy the server when done:**

```bash
terraform destroy
```

## Manual Steps (Alternative to run-example.sh)

If you prefer to run commands manually:

```bash
cd full-example
terraform init
terraform plan     # See what will be created
terraform apply     # Create the server
terraform destroy   # Clean up when done
```