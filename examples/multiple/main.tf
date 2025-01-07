

terraform {
  required_providers {
    infrahub = {
      source  = "registry.marcomartinez.ch/marcom4rtinez/infrahub"
      version = "1.0"
    }
  }
}

provider "infrahub" {
  api_key         = "180e8659-2f40-400d-36ac-c513d60a378c"
  infrahub_server = "http://localhost:8000"
  branch          = "main"
}

data "infrahub_artifact" "name" {
  artifact_id     = ""
  infrahub_server = "http://localhost:8000"
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

# output "device_res" {
#   value = infrahub_device.device_res.edges_node_topology_node_name_value
# }

# data "infrahub_devicequery" "tset" {
#   device_name = infrahub_device.device_res.edges_node_name_value
# }

# output "device_name_queried" {
#   value = data.infrahub_devicequery.tset.edges_node_role_value
# }

# data "infrahub_country" "germany" {
#   country_name = "Germany"
# }

# data "infrahub_topology" "de1-pod1" {
#   topology_name = "de1-pod1"
# }

# data "infrahub_devicetype" "ccs" {
#   device_type_name = "CCS-720DP-48S-2F"
# }

# data "infrahub_autonomoussystem" "AS174" {
#   as_name = "AS174"
# }

# data "infrahub_platform" "Arista" {
#   platform_name = "Arista EOS"
# }

# data "infrahub_ipaddress" "mgmt_address" {
#   ip_address_value = "10.0.0.1/24"
# }

# resource "infrahub_device" "device_res" {
#   name_value              = "switch27"
#   asn_node_id             = data.infrahub_autonomoussystem.AS174.id
#   device_type_node_id     = data.infrahub_devicetype.ccs.id
#   location_node_id        = data.infrahub_country.germany.id
#   platform_node_id        = data.infrahub_platform.Arista.id
#   primary_address_node_id = data.infrahub_ipaddress.mgmt_address.id
#   status_value            = "active"
#   topology_node_id        = data.infrahub_topology.de1-pod1.id
#   role_value              = "leaf"
# }

