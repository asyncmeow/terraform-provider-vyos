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

data "vyos_interfaces_ethernet" "eth0" {
  name = "eth2"
}

output "test" {
  value = data.vyos_interfaces_ethernet.eth0
}