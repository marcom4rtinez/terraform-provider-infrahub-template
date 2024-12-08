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

# resource "infrahub_device" "device_res" {
#   device_name = "fra05-pod6-leaf4"
#   # name        = "fra05-pod1-leaf1"
#   role = "leaf"
#   # id          = "1802e1f2-bc07-e55b-395f-c515fdfc0604"
# }

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

resource "infrahub_device" "device_res" {
  device_name                        = "test_device"
  edges_node_asn_node_id             = "180d513c-c700-f27f-36ae-c5147a57daa5"
  edges_node_device_type_node_id     = "180d513e-0bc4-6766-36af-c514f3173dcf"
  edges_node_location_node_id        = "180d5144-a48a-161d-36ad-c51efb37fd26"
  edges_node_platform_node_id        = "180d513d-b76a-dd27-36a0-c51c912c8d09"
  edges_node_primary_address_node_id = "180d52a5-84a6-a85d-36ac-c511c32e48bf"
  edges_node_status_value            = "active"
  edges_node_topology_node_id        = "180d514e-d1c0-61df-36af-c51a2a0d705b"
  edges_node_role_value              = "client"
  # edges_node_description_id           = ""
  # edges_node_description_value        = ""
  # edges_node_id                       = ""
  # edges_node_name_value               = ""
  # edges_node_role_id                  = ""
  # edges_node_status_id                = ""
  # edges_node_topology_node_name_value = ""
}



# resource "infrahub_device" "device_res" {
#   device_name = "fra05-pod6-leaf4"
#   # name        = "fra05-pod1-leaf1"
#   role = "leaf"
#   # id          = "1802e1f2-bc07-e55b-395f-c515fdfc0604"
# }
