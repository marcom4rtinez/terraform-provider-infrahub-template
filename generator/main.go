// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateTerraformDataSource(parsedQuery *InputGraphQLQuery) (string, error) {
	structName := parsedQuery.QueryName + "DataSource"
	data := DataSourceTemplateData{
		QueryName:       parsedQuery.QueryName,
		ObjectName:      parsedQuery.ObjectName,
		Required:        parsedQuery.Required,
		StructName:      structName,
		Fields:          parsedQuery.Fields,
		GenqlientFields: parsedQuery.GenqlientFields,
	}

	// Render the template
	caser := cases.Title(language.English)
	datasourceTemplate, err := template.New("datasource").Funcs(template.FuncMap{
		"title": caser.String,
	}).Parse(string(datasourceTemplateContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = datasourceTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateTerraformResource(parsedQuery *InputGraphQLQuery) (string, error) {
	structName := parsedQuery.QueryName + "Resource"
	data := ResourceTemplateData{
		QueryName:               parsedQuery.QueryName,
		ObjectName:              parsedQuery.ObjectName,
		Required:                parsedQuery.Required,
		StructName:              structName,
		Fields:                  parsedQuery.Fields,
		GenqlientFields:         parsedQuery.GenqlientFields,
		GenqlientFieldsModify:   parsedQuery.genqlientFieldsModify,
		GenqlientFieldsReadOnly: parsedQuery.genqlientFieldsReadOnly,
	}

	// Render the template
	caser := cases.Title(language.English)
	resourceTemplate, err := template.New("resource").Funcs(template.FuncMap{
		"title": caser.String,
	}).Parse(string(resourceTemplateContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = resourceTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateTerraformProvider(components TerraformComponents) (string, error) {
	data := ProviderSourceTemplateData{
		DataSources: components.dataSources,
		Resources:   components.resources,
	}

	// Render the template
	caser := cases.Title(language.English)
	providerTemplate, err := template.New("provider").Funcs(template.FuncMap{
		"title": caser.String,
	}).Parse(string(providerTemplateContent))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = providerTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func readAndGenerateDataSourcesAndResources(graphqlQuery string) (string, string, error) {

	parsedQuery, err := ParseGraphQLQuery(graphqlQuery)

	if err != nil {
		fmt.Println("Error parsing GraphQL query:", err)
		os.Exit(1)
	}

	if parsedQuery.ResourceType == DataSource {
		code, err := GenerateTerraformDataSource(parsedQuery)
		if err != nil {
			fmt.Println("Error generating Terraform data source:", err)
			os.Exit(1)
		}
		file, err := os.Create(fmt.Sprintf("../internal/provider/%s_data_source.go", parsedQuery.QueryName))
		if err != nil {
			fmt.Println("Error creating the file:", err)
			return "", "", err
		}
		defer file.Close()

		_, err = file.WriteString(code)
		if err != nil {
			fmt.Println("Error writing to the file:", err)
			return "", "", err
		}

		fmt.Printf("Content written to %s_data_source.go file successfully!\n", parsedQuery.QueryName)
		return parsedQuery.QueryName, "", nil
	} else if parsedQuery.ResourceType == Resource {
		code, err := GenerateTerraformResource(parsedQuery)
		if err != nil {
			return "", "", fmt.Errorf("Error generating Terraform resource: %s", err)
		}
		file, err := os.Create(fmt.Sprintf("../internal/provider/%s_resource.go", parsedQuery.QueryName))
		if err != nil {
			return "", "", fmt.Errorf("Error creating the file: %s", err)
		}
		defer file.Close()

		_, err = file.WriteString(code)
		if err != nil {
			return "", "", fmt.Errorf("Error writing to the file: %s", err)
		}

		fmt.Printf("Content written to %s_resource.go file successfully!\n", parsedQuery.QueryName)
		return "", parsedQuery.QueryName, nil
	}

	return "", "", fmt.Errorf("No Resource or DataSource")

}

func main() {
	gqlDir := "gql"

	var dataSources, resources []string

	err := filepath.Walk(gqlDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if filepath.Ext(path) == ".gql" {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				dataSourceName, resourceName, err := readAndGenerateDataSourcesAndResources(string(data))
				if err == nil {
					if dataSourceName != "" {
						dataSources = append(dataSources, dataSourceName)
					} else if resourceName != "" {
						resources = append(resources, resourceName)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	readAndGenerateProvider(
		TerraformComponents{
			dataSources: dataSources,
			resources:   resources,
		})
}

func readAndGenerateProvider(components TerraformComponents) {

	code, err := GenerateTerraformProvider(components)

	if err != nil {
		return
	}

	file, err := os.Create("../internal/provider/provider.go")
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}

	fmt.Printf("Content written to provider.go file successfully!\n")
}
