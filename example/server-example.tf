
# Example: Provision a web server
resource "ics_bare_metal_server" "web_server" {
  instance_type    = "c2.small"
  location         = "FRA1"
  operating_system = "Ubuntu 24.04"
  hostname         = "web-server"
  friendly_name    = "ICS Terraform Demo Web Server"
  ssh_key_labels   = [ics_ssh_key.example.label]
}

# Example: Provision a database server
resource "ics_bare_metal_server" "db_server" {
  instance_type    = "c2.small"
  location         = "FRA1"
  operating_system = "Ubuntu 24.04"
  hostname         = "db-server"
  friendly_name    = "ICS Terraform Demo DB Server"
  ssh_key_labels   = [ics_ssh_key.example.label]
}

# Output server information
output "server_info" {
  description = "Information about provisioned servers"
  value = {
    web_server = {
      service_id = ics_bare_metal_server.web_server.service_id
      public_ip  = ics_bare_metal_server.web_server.public_ip
      hostname   = ics_bare_metal_server.web_server.hostname
    }
    db_server = {
      service_id = ics_bare_metal_server.db_server.service_id
      public_ip  = ics_bare_metal_server.db_server.public_ip
      hostname   = ics_bare_metal_server.db_server.hostname
    }
  }
}