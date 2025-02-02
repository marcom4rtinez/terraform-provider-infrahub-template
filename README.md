# Terraform Provider Infrahub (Terraform Plugin Framework)
This repository is intended to be used as Template. To get your provider up and running do the following steps.


## Create your own Repository
Go to the button at the top of the GitHub Page that says `Use this template` then further `Create new repository`. Choose an adequate name like terraform-provider-infrahub.

> [!CAUTION]
> The name of the repository has to start with `terraform-provider-` that is given by terraform
> The provider name will be the prefix of your provider e.g. `terraform-provider-infrahub` will become `terraform-provider-infrahub-main`

Once Repository is created go ahead and clone it.

## Generate the provider
In the cloned repository copy all your gql queries you wish to be processed into `gql/`.

Change the below value to your Infrahub Instance
```
export INFRAHUB_SERVER="http://localhost:8000"
```

Now generate the SDK, provider and documentation
```
make all
```

The Provider is now fully functional. Either move to the section local development to test it or move on to the deployment part of this guide.

## Prerequisites for deployment
Make sure you set these values to be able to package and sign the provider.

1. Set ENV GPG_FINGERPRINT `export GPG_FINGERPRINT=9A52F2BE41E9C446A902C723B53E44105C84C057`
2. Set ENV GPG_PUBLIC_KEY `export GPG_PUBLIC_KEY=$(gpg --armor --export $GPG_FINGERPRINT)`
3. Set ENV GITHUB_TOKEN `export GITHUB_TOKEN=XXXXX`
4. Set ENV TERRAFORM_REGISTRY `export TERRAFORM_REGISTRY_ENDPOINT="http://localhost:8080/v1/providers/marcom4rtinez/infrahub-main/upload"`
5. Set ENV RELEASE_URL `export RELEASE_URL="https://github.com/marcom4rtinez/terraform-provider-infrahub/releases/download"`

## Deploy to a private Registry
This guide assumes you are using the Registry from https://github.com/marcom4rtinez/terraform-registry. Otherwise adjust `GNUmakefile` to suite your needs.

In this last step the provider will be deployed to the Registry and stored in GitHub Releases. If not done yet, please make a git tag
```
git tag -a v1.4 -m "my awesome Infrahub Terraform Provider"
```

Next execute the deployment process.
```
make generate_deploy
```


# Miscellaneous

## Generate automatic documentation

Documentation is generated automatically as part of the building process, but it can be done manually using `make generate`, there are possibilities to add more examples. Examples can be added in `example/` consult `example/README.md` for more information. All documentation is available in `docs/`.


## Prerequisites for local development
To execute custom terraform provider locally set the variable in `$USER/.terraformrc`
```bash
provider_installation {

  dev_overrides {
      "registry.terraform.io/marcom4rtinez/infrahub" = "/Users/marco/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

Build custom provider with SDK and generator.
```
# Builds SDK from gql/, generates files for provider
# May overwrite changes in internal/provider
make all
```

To debug the code generated code / add more
```
# Installs local provider to local dev environment
# Will just execute files as in internal/provider
make
```

Once installed with either `make` or `make all` it can be used in `main.tf`
```
# Make sure overwrite matches .terraformrc

terraform {
  required_providers {
    infrahub = {
      source  = "registry.marcomartinez.ch/marcom4rtinez/infrahub"
      version = "1.0"
    }
  }
}
```