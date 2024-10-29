terraform {
  required_providers {
    infrahub = {
      source  = "registry.terraform.io/marcom4rtinez/infrahub"
      version = "1.0"
    }
  }
}

resource "infrahub_example" "example" {
  configurable_attribute = "some-value"
}

