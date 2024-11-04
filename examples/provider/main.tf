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
