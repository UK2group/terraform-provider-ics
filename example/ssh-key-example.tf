# Example: Create an SSH key for server access
resource "ics_ssh_key" "example" {
  label      = "my-terraform-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDL8/Example+Public+Key+Content+Here+Replace+With+Your+Own+Public+Key+AAAAB3NzaC1yc2EAAAADAQABAAABgQDL8 user@hostname"
}

# Output the SSH key details
output "ssh_key_info" {
  description = "Information about the created SSH key"
  value = {
    id          = ics_ssh_key.example.id
    label       = ics_ssh_key.example.label
  }
}