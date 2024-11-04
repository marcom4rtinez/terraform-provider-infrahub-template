terraform {
  required_providers {
    infrahub = {
      source  = "registry.terraform.io/marcom4rtinez/infrahub"
      version = "1.0"
    }
  }
}

provider "infrahub" {
  api_key         = "XXX"
  infrahub_server = "10.0.0.1"
}

data "infrahub_device" "fra05-pod1-leaf1" {
  device_name = "fra05-pod1-leaf1"
}

output "device_name_output" {
  value = data.infrahub_device.fra05-pod1-leaf1.name
}

output "device_id_output" {
  value = data.infrahub_device.fra05-pod1-leaf1.id
}

output "device_role_output" {
  value = data.infrahub_device.fra05-pod1-leaf1.role
}
