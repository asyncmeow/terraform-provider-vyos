terraform {
  required_providers {
    vyos = {
      source = "registry.terraform.io/asyncmeow/vyos"
    }
  }
}

provider "vyos" {
  host     = "https://172.26.0.57"
  key      = "test"
  insecure = true
}

data "vyos_ethernet_interface" "eth0" {
  name = "eth2"
}

output "test" {
  value = data.vyos_ethernet_interface.eth0
}