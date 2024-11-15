resource "infrahub_node" "my_device" {
  kind = "NetworkDevice"
  attributes = {
    "name" = "lon-edg01"
  }
}