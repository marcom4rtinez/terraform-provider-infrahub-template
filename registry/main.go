package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"
	"text/template"
)

// Structures to match the JSON template

type TemplateData struct {
	Version        string
	Protocols      []string
	Platforms      []Platform
	ProviderName   string
	BaseURL        string
	KeyID          string
	AsciiArmor     string
	TrustSignature string
	Source         string
	SourceURL      string
}

type Platform struct {
	OS           string
	Arch         string
	PlatformName string
}

func main() {
	// Define flags for CLI inputs
	version := flag.String("version", "", "Version of the Provider")
	protocols := flag.String("protocols", "", "Comma-separated list of protocols")
	os := flag.String("os", "", "Operating system")
	arch := flag.String("arch", "", "Architecture")
	providerName := flag.String("provider_name", "", "Name of the provider")
	baseURL := flag.String("base_url", "", "Download URL for the artifact")
	keyID := flag.String("key_id", "", "GPG key ID")
	asciiArmor := flag.String("ascii_armor", "", "GPG public key in ASCII armor format")
	source := flag.String("source", "", "Source of the GPG key")
	sourceURL := flag.String("source_url", "", "Source URL of the GPG key")

	// Parse flags
	flag.Parse()

	// Create Platform combinations
	platforms := []Platform{}
	for _, current_os := range strings.Split(*os, ",") {
		for _, current_arch := range strings.Split(*arch, ",") {
			current_platform := Platform{
				OS:           current_os,
				Arch:         current_arch,
				PlatformName: fmt.Sprintf("%s_%s", current_os, current_arch),
			}
			platforms = append(platforms, current_platform)
		}
	}

	// Create the data structure to pass to the template
	templateData := TemplateData{
		Version:        *version,
		Protocols:      strings.Split(*protocols, ","),
		Platforms:      platforms,
		ProviderName:   *providerName,
		BaseURL:        *baseURL,
		KeyID:          *keyID,
		AsciiArmor:     strings.ReplaceAll(strings.ReplaceAll(*asciiArmor, "\n", "\\n"), "-----BEGIN PGP PUBLIC KEY BLOCK-----", "-----BEGIN PGP PUBLIC KEY BLOCK-----\\nVersion: GnuPG v1"),
		TrustSignature: "",
		Source:         *source,
		SourceURL:      *sourceURL,
	}

	tmpl, err := template.New("registryTemplate").Parse(registryTemplateContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
	fmt.Println(buf.String())
}
