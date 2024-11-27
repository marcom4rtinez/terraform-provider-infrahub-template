// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type InputGraphQLQuery struct {
	QueryName       string
	ObjectName      string
	Required        string
	Fields          []Field
	GenqlientFields []GenqlientField
}

type Field struct {
	Name string
	Type string
}

type GenqlientField struct {
	Name  string
	Query string
}

type DataSourceTemplateData struct {
	QueryName       string
	ObjectName      string
	Required        string
	StructName      string
	Fields          []Field
	GenqlientFields []GenqlientField
}
type ProviderSourceTemplateData struct {
	DataSources []string
	Resources   []string
	Functions   []string
}

func GraphQLToTerraform(graphqlType string) string {
	switch graphqlType {
	case "String":
		return "types.String"
	case "Int":
		return "types.Int64"
	case "Float":
		return "types.Float64"
	case "Boolean":
		return "types.Bool"
	default:
		return "types.String"
	}
}

func ParseGraphQLQuery(query string) (*InputGraphQLQuery, error) {
	lines := strings.Split(query, "\n")

	var queryName, required, parentPrefix, objectName string
	var fields []Field
	var inBlock bool
	var prefixList, prefixListImmutable []string

	for number, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "query ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				containsBracket := strings.IndexByte(parts[1], '(')
				if containsBracket != -1 {
					queryName = parts[1][:containsBracket]
				} else {
					queryName = parts[1]
				}
				queryName = strings.ToLower(string(queryName[0])) + queryName[1:]
			}
		} else if number == 1 {
			// } else if strings.Contains(line, ": $") {
			// This identifies the required field (e.g., name__value: $device_name)
			if strings.Contains(line, ":") {
				parts := strings.Split(line, ":")
				required = parts[1][strings.Index(parts[1], "$")+1 : strings.Index(parts[1][strings.Index(parts[1], "$"):], " ")+strings.Index(parts[1], "$")]
				required = strings.TrimRight(required, ")")
				objectNameParts := strings.Split(parts[0], "(")
				objectName = objectNameParts[0]
			} else {
				parts := strings.Split(line, " ")
				objectName = parts[0]
			}
		} else if strings.HasSuffix(line, " {") {
			inBlock = true
			prefix := line[:len(line)-2]
			prefixList = append(prefixList, prefix)
			if strings.Contains(prefix, "_") {
				prefixListImmutable = append(prefixListImmutable, prefix)
			}
			parentPrefix = parentPrefix + prefix + "_"
		} else if line == "}" {
			inBlock = false
			if strings.Count(parentPrefix, "_") < 2 {
				parentPrefix = ""
				break
			}
			// remove last _ and length of last prefix added, workaround for underscores in schema
			parentPrefix = parentPrefix[:len(parentPrefix)-1-len(prefixList[len(prefixList)-1])]
			prefixList = prefixList[:len(prefixList)-1]
			// parentPrefix = parentPrefix[:strings.LastIndex(parentPrefix[:strings.LastIndex(parentPrefix, "_")], "_")+1]
		} else if inBlock {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				fields = append(fields, Field{
					Name: parentPrefix + strings.TrimSpace(parts[0]),
					Type: "String",
				})
				if strings.Contains(parts[0], "_") {
					prefixListImmutable = append(prefixListImmutable, parts[0])
				}
			}
		}
	}

	customSplit := func(str string, exceptions []string) []string {
		var result []string
		var currentWord string

		for _, char := range str {
			if char == '_' {
				isException := false
				for _, exception := range exceptions {
					if strings.HasPrefix(exception, currentWord) {
						if len(currentWord) == len(exception) {
							break
						}
						isException = true
						break
					}
				}
				if !isException {
					result = append(result, currentWord)
					currentWord = ""
				} else {
					currentWord += string(char)
				}
			} else {
				currentWord += string(char)
			}
		}
		result = append(result, currentWord)
		return result
	}

	var genqlientFields []GenqlientField

	for _, entry := range fields {
		parts := customSplit(entry.Name, prefixListImmutable)

		// Capitalize each part except for the first one
		caser := cases.Title(language.English)
		for i := range parts {
			// Capitalize the first letter of each part
			parts[i] = caser.String(parts[i])
			if required != "" {
				if parts[i] == "Edges" {
					parts[i] = "Edges[0]"
				}
			} else {
				if parts[i] == "Edges" {
					parts[i] = "Edges[i]"
				}
			}
		}

		// Join the parts using a dot separator
		genqlientFields = append(genqlientFields, GenqlientField{
			Name:  entry.Name,
			Query: objectName + "." + strings.Join(parts, "."),
		})
	}

	if queryName == "" {
		return nil, fmt.Errorf("failed to parse GraphQL query: missing query name")
	}

	return &InputGraphQLQuery{
		QueryName:       queryName,
		ObjectName:      objectName,
		Required:        required,
		Fields:          fields,
		GenqlientFields: genqlientFields,
	}, nil
}

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

func GenerateTerraformProvider(dataSourcesList []string) (string, error) {
	data := ProviderSourceTemplateData{
		DataSources: dataSourcesList,
	}

	// Render the template
	caser := cases.Title(language.English)
	datasourceTemplate, err := template.New("provider").Funcs(template.FuncMap{
		"title": caser.String,
	}).Parse(string(providerTemplateContent))
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

func readAndGenerateDataSources(graphqlQuery string) (string, error) {

	// Parse the query
	parsedQuery, err := ParseGraphQLQuery(graphqlQuery)
	if err != nil {
		fmt.Println("Error parsing GraphQL query:", err)
		os.Exit(1)
	}

	// Generate the Terraform data source code
	code, err := GenerateTerraformDataSource(parsedQuery)
	if err != nil {
		fmt.Println("Error generating Terraform data source:", err)
		os.Exit(1)
	}

	file, err := os.Create(fmt.Sprintf("../internal/provider/%s_data_source.go", parsedQuery.QueryName))
	if err != nil {
		fmt.Println("Error creating the file:", err)
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(code)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return "", err
	}

	fmt.Printf("Content written to %s_data_source.go file successfully!\n", parsedQuery.QueryName)
	return parsedQuery.QueryName, nil
}

func main() {
	gqlDir := "gql"

	var dataSources []string

	err := filepath.Walk(gqlDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			dataSourceName, err := readAndGenerateDataSources(string(data))
			if err == nil {
				dataSources = append(dataSources, dataSourceName)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	readAndGenerateProvider(dataSources)
}

func readAndGenerateProvider(dataSources []string) {

	code, err := GenerateTerraformProvider(dataSources)

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
