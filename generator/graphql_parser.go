// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GraphQLToTerraformTypes(graphqlType string) string {
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
	var resourceType ResourceType
	var result InputGraphQLQuery
	var err error
	lines := strings.Split(query, "\n")

	for _, line := range lines {
		if strings.Contains(line, "mutation") {
			resourceType = Resource
			break
		} else if strings.Contains(line, "query") {
			resourceType = DataSource
			break
		}
	}

	if resourceType == DataSource {
		result, err = parseDataSourceInput(lines)
		result.ResourceType = DataSource

	} else if resourceType == Resource {
		result, err = parseResourceInput(lines)
		result.ResourceType = Resource
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func parseResourceInput(lines []string) (InputGraphQLQuery, error) {
	var queryName, required, objectName, parentPrefix string
	var inBlock bool
	var prefixList, prefixListImmutable []string
	var fields []Field
	var genqlientFields, genqlientFieldsModify, genqlientFieldsReadOnly []GenqlientField

	index := 0
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
			index = number
		} else if index != 0 && number == index+1 {
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
		}
	}

	for number, line := range lines[index:] {

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "query ") || number == 1 {
			//already handled
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

	for _, entry := range fields {
		parts := customSplit(entry.Name, prefixListImmutable)

		// Capitalize each part except for the first one
		caser := cases.Title(language.English)
		var filtered, noPrefix, plain []string

		for i := range parts {
			// Capitalize the first letter of each part
			parts[i] = caser.String(parts[i])
			plain = append(plain, parts[i])
			if parts[i] != "Edges" && parts[i] != "Node" {
				noPrefix = append(noPrefix, parts[i])
				filtered = append(filtered, parts[i])
			}
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
		// skippedEdges = false
		// skippedNode = false

		for _, x := range [][]string{parts, noPrefix, plain} {
			if len(x) > 0 && x[len(x)-1] == "Id" {
				x[len(x)-1] = "GetId()"
			}
		}

		newField := GenqlientField{
			Name:                   entry.Name,
			Query:                  objectName + "." + strings.Join(parts, "."),
			QueryNoPrefixReplaceId: strings.Join(noPrefix, "."),
			InputObjectNames:       strings.Join(filtered, "."),
			PlainObject:            strings.Join(plain[2:], "."),
		}

		if strings.Count(strings.ToLower(newField.Query), "node") < 2 && strings.Count(strings.ToLower(newField.Query), "id") < 1 {
			genqlientFieldsModify = append(genqlientFieldsModify, newField)
		} else if strings.Count(strings.ToLower(newField.Query), "node") >= 2 && strings.Count(strings.ToLower(newField.Query), "id") >= 1 {
			genqlientFieldsModify = append(genqlientFieldsModify, newField)
		} else {
			genqlientFieldsReadOnly = append(genqlientFieldsReadOnly, newField)
		}
		genqlientFields = append(genqlientFields, newField)
	}

	if queryName == "" {
		return InputGraphQLQuery{}, fmt.Errorf("failed to parse GraphQL query: missing query name")
	}

	return InputGraphQLQuery{
		QueryName:               queryName,
		ObjectName:              objectName,
		Required:                required,
		Fields:                  fields,
		GenqlientFields:         genqlientFields,
		genqlientFieldsReadOnly: genqlientFieldsReadOnly,
		genqlientFieldsModify:   genqlientFieldsModify,
	}, nil
}

func parseDataSourceInput(lines []string) (InputGraphQLQuery, error) {
	var queryName, required, objectName, parentPrefix string
	var fields []Field
	var genqlientFields []GenqlientField
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
		return InputGraphQLQuery{}, fmt.Errorf("failed to parse GraphQL query: missing query name")
	}

	return InputGraphQLQuery{
		QueryName:       queryName,
		ObjectName:      objectName,
		Required:        required,
		Fields:          fields,
		GenqlientFields: genqlientFields,
	}, nil
}
