terraform {
  required_providers {
    infrahub = {
      source  = "registry.terraform.io/marcom4rtinez/infrahub"
      version = "1.0"
    }
  }
}

provider "infrahub" {
  api_key         = "180e8659-2f40-400d-36ac-c513d60a378c"
  infrahub_server = "localhost"
}

# data "infrahub_devices" "example" {
# }

# output "devices_example" {
#   value = data.infrahub_devices.example
# }

# output "device_id_output" {
#   value = data.infrahub_device.fra05-pod1-leaf1.id
# }

# output "device_role_output" {
#   value = data.infrahub_device.fra05-pod1-leaf1.role
# }

resource "infrahub_device" "device_res" {
  device_name = "fra05-pod6-leaf4"
  # name        = "fra05-pod1-leaf1"
  role = "leaf"
  # id          = "1802e1f2-bc07-e55b-395f-c515fdfc0604"
}

# output "device_resu" {
#   value = infrahub_device.device_res
# }

# data "infrahub_device" "fra05-pod1-leaf1" {
#   device_name = "fra05-pod1-leaf1"
# }

# output "device_name_output" {
#   value = data.infrahub_device.fra05-pod1-leaf1
# }

# data "infrahub_interface" "ethernet12" {
#   interface_name = "ge-0/0/0"
# }

# output "ethernet1_output" {
#   value = data.infrahub_interface.ethernet12
# }

# data "infrahub_devices" "all_devices" {
# }

# output "all_devices_output" {
#   value = data.infrahub_devices.all_devices
# }

# data "infrahub_accounts" "all_accounts" {
# }

# output "all_accounts_output" {
#   value = data.infrahub_accounts.all_accounts
# }

# data "infrahub_bgpsessions" "all_bgp_sessions" {
# }

# output "all_bgp_sessions_output" {
#   value = data.infrahub_bgpsessions.all_bgp_sessions
# }
