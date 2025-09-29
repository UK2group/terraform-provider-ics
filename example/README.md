# ICS Terraform Provider Examples

This directory contains example configurations for the ICS Terraform Provider.

## Examples

### SSH Key Management (`ssh-key-example.tf`)
Demonstrates how to create and manage SSH keys using the `ics_ssh_key` resource.

### Server Provisioning (`server-example.tf`)
Shows how to provision bare metal servers with SSH key access using the `ics_bare_metal_server` resource.

### Complete Example (`main.tf`)
A comprehensive example that combines SSH key creation and server provisioning.

## Prerequisites

1. **ICS API Token**: Obtain your API token from the ICS dashboard
2. **SSH Public Key**: Have your SSH public key ready for server access

## Setup

1. **Set your API token** (recommended method):
   ```bash
   export ICS_API_TOKEN="your-actual-api-token-here"
   ```

2. **Replace the example SSH key** in the configuration files with your actual public key:
   ```bash
   # Replace this placeholder:
   ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDL8/Example+Public+Key+Content...

   # With your actual public key:
   ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@hostname
   ```

### 2. Navigate to the Example Directory

```bash
cd full-example
```

### 3. Initialize Terraform

```bash
terraform init
```

This will download and configure the locally installed ICS provider.

### 4. Plan the Deployment

First, see what Terraform will do:

```bash
terraform plan
```

This will show you:
- Available server configurations
- Which server type and location will be used
- Available operating systems
- The planned server configuration

### 5. Apply the Configuration

Deploy the server:

```bash
terraform apply
```

**Important**: This will create a real bare metal server and you will be charged for it!

## Expected Outputs

### Immediate Outputs (after `terraform apply` starts)

- **Service ID**: Available immediately after the order is placed
- **Available Configurations**: Shows all server types with auto-provision inventory
- **Available Operating Systems**: Shows OS options for your server type/location

### Post-Provisioning Outputs (after provisioning completes)

- **Public IP Address**: The server's public IP (available after 5-30 minutes)
- **Server Details**: Complete information about the provisioned server
- **MAC Address**: Network interface MAC address
- **Datacenter Info**: Location and vendor details

## Example Output

```bash
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

server_details = {
  "datacenter_name" = "Dallas"
  "friendly_name" = "ICS Terraform Demo Server"
  "hostname" = "ics-terraform-demo"
  "instance_type" = "c1.small"
  "location" = "NYC1"
  "operating_system" = "Ubuntu 24.04"
  "public_ip" = "123.2.123.132"
  "service_id" = 12345
}
```

## Error Handling

If you specify invalid values, the provider gives helpful error messages:

### Invalid Instance Type + Location Combination
```bash
Error: Invalid Instance Type and Location Combination

Instance type 'c1.xlarge' is not available in location 'NYC1'

Available instance types in location 'NYC1': [c1.small c1.medium]

Available locations for instance type 'c1.xlarge': [SLC1 ORD1]

All available combinations with inventory:
  c1.small: [NYC1 SLC1]
  c1.medium: [NYC1 ORD1 SLC1]
  c1.xlarge: [SLC1 ORD1]
```

### Invalid Operating System
```bash
Error: Invalid Operating System

Operating system 'Windows Server 2022' is not available for instance type 'c1.small' in location 'NYC1'

Available operating systems: [Ubuntu 24.04 Ubuntu 22.04 Debian 12 Debian 11 CentOS 8]
```

## Monitoring Provisioning Progress

### Option 1: Wait for Terraform
Terraform will automatically wait up to 30 minutes for provisioning to complete. You'll see logs like:

```
ics_bare_metal_server.example: Creating...
ics_bare_metal_server.example: Still creating... [30s elapsed]
ics_bare_metal_server.example: Still creating... [1m0s elapsed]
...
ics_bare_metal_server.example: Creation complete after 8m42s [service_id=12345]
```

### Option 2: Check Status Manually
You can check the current status without waiting:

```bash
terraform refresh
terraform output server_public_ip
```

### Option 3: ICS Portal
Log into the ICS portal and check the server status under your services.

## Cleanup

⚠️ **Important**: Don't forget to destroy the server when you're done to avoid ongoing charges!

```bash
terraform destroy
```

This will automatically cancel the hourly-billed server via the API.

## Customization

You can customize the server by modifying the `ics_bare_metal_server` resource in `main.tf`:

```hcl
resource "ics_bare_metal_server" "example" {
  instance_type    = "c1.medium"     # Choose different instance type
  location         = "SLC11"           # Choose different location
  operating_system = "Debian 12"      # Choose different OS
  hostname         = "my-server"      # Customize hostname
  domain           = "mydomain.com"   # Customize domain
  friendly_name    = "My Server"      # Customize friendly name
  notes           = "My custom notes" # Add custom notes
}
```

The provider will automatically validate that your chosen combination is available and provide helpful error messages if not.

## Troubleshooting

### Provider Not Found
If you get an error about the provider not being found, rebuild and reinstall:

```bash
cd .. # Go back to provider root
make install
cd full-example
rm -rf .terraform .terraform.lock.hcl  # Clean up
terraform init
```

### First Time Setup
If this is your first time running the example, make sure the provider is built and installed:

```bash
cd .. # Go back to provider root
make install  # This builds and installs the provider locally
cd full-example
terraform init
```

### API Authentication Errors
Make sure your `ICS_API_TOKEN` environment variable is set correctly:

```bash
echo $ICS_API_TOKEN  # Should show your token
```

### No Auto-Provision Inventory
If no servers are available for auto-provisioning, the example will fail. Check the `available_server_configs` output to see what's available, or try again later.

### Provisioning Timeout
If provisioning takes longer than 30 minutes, you can:
1. Increase the timeout in the provider code
2. Check the ICS portal for server status
3. Contact ICS support if there are infrastructure issues

## Support

For issues with:
- **The Terraform Provider**: Check the provider code and documentation
- **ICS API**: Contact Ingenuity Cloud Services support
- **Server Provisioning**: Contact ICS support with your service ID