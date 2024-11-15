terraform {
  required_providers {
    infrahub = {
      source  = "registry.terraform.io/marcom4rtinez/infrahub"
      version = "1.0"
    }
  }
}

provider "infrahub" {
  api_key = "not_needed_for_data_sources"
  //hardcoded localhost for artifact query, otherwise using this value
  infrahub_server = "localhost"
}

data "infrahub_artifact" "artifact1" {
  artifact_id = "1808091d-23dc-9424-3956-c516ad3482e9"
}


output "artifact1" {
  value = data.infrahub_artifact.artifact1.content
}
