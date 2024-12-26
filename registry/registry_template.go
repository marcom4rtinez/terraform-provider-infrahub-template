package main

const registryTemplateContent = `{
"versions": [
	{
	"version": "{{.Version}}",
	"protocols": [{{range $index, $protocol := .Protocols}}{{if $index}}, {{end}}"{{$protocol}}"{{end}}],
	"platforms": [
			{{- $BaseURL :=  .BaseURL  }}
			{{- $ProviderName :=  .ProviderName  }}
			{{- $Version :=  .Version  }}
			{{- $KeyID :=  .KeyID  }}
			{{- $AsciiArmor :=  .AsciiArmor  }}
			{{- $TrustSignature :=  .TrustSignature  }}
			{{- $Source :=  .Source  }}
			{{- $SourceURL :=  .SourceURL  }}
			{{- range $index, $platform := .Platforms }}
			{{- if and ($index) (gt $index 0)}},{{end}}
				{
					"os": "{{.OS}}",
					"arch": "{{.Arch}}",
					"filename": "{{$ProviderName}}_{{$Version}}_{{.OS}}_{{.Arch}}.zip",
					"download_url": "{{$BaseURL}}/v{{$Version}}/{{$ProviderName}}_{{$Version}}_{{.OS}}_{{.Arch}}.zip",
					"shasums_url": "{{$BaseURL}}/v{{$Version}}/{{$ProviderName}}_{{$Version}}_SHA256SUMS",
					"shasums_signature_url": "{{$BaseURL}}/v{{$Version}}/{{$ProviderName}}_{{$Version}}_SHA256SUMS.sig",
					"shasum": "{{$ProviderName}}_{{$Version}}_{{.OS}}_{{.Arch}}.zip_shasum",
					"signing_keys": {
						"gpg_public_keys": [
							{
								"key_id": "{{$KeyID}}",
								"ascii_armor": "{{$AsciiArmor}}",
								"trust_signature": "{{$TrustSignature}}",
								"source": "{{$Source}}",
								"source_url": "{{$SourceURL}}"
							}
						]
					}
				}
			{{- end}}
			]
		}
	]
}`
